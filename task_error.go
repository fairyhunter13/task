package task

import (
	"sync"
	"sync/atomic"

	"github.com/panjf2000/ants"
)

const (
	channelClosed uint64 = 1
	channelOpen   uint64 = 0
	// DefaultBufferSize specifies the default size of channel.
	DefaultBufferSize = 5
)

// ErrorManager contains Manager but especially for error function.
type ErrorManager struct {
	chErr    chan error
	isClosed uint64
	bufSize  int
	wg       *sync.WaitGroup
}

// NewErrorManager initialize the new error manager.
func NewErrorManager(bufSize int) *ErrorManager {
	em := &ErrorManager{
		bufSize: bufSize,
	}
	em.init()
	return em
}

func (em *ErrorManager) init() {
	if em.bufSize < 1 {
		em.bufSize = DefaultBufferSize
	}
	if em.chErr == nil || em.isChannelClosed() {
		em.chErr = make(chan error, em.bufSize)
		atomic.StoreUint64(&em.isClosed, channelOpen)
	}
	if em.wg == nil {
		em.wg = new(sync.WaitGroup)
	}
}

func (em *ErrorManager) isChannelClosed() bool {
	return atomic.LoadUint64(&em.isClosed) == channelClosed
}

// ClosureErr defines closure that returns error.
type ClosureErr func() (err error)

// Run runs the closure error function.
func (em *ErrorManager) Run(fn ClosureErr) {
	if fn == nil {
		return
	}
	em.init()
	em.wg.Add(1)
	ants.Submit(func() {
		defer em.wg.Done()
		em.chErr <- fn()
	})
}

// ErrChan returns the receiving error channel of this error manager.
func (em *ErrorManager) ErrChan() <-chan error {
	em.init()
	return em.chErr
}

// WaitClose wait all go routines to complete and close the channel in the separate go routine.
func (em *ErrorManager) WaitClose() {
	em.init()
	ants.Submit(func() {
		em.wg.Wait()
		em.close()
	})
}

// Error returns the first error from the Run execution of the fn closure.
func (em *ErrorManager) Error() (err error) {
	em.WaitClose()

	for errTemp := range em.chErr {
		if err != nil {
			continue
		}

		err = errTemp
	}
	return
}

func (em *ErrorManager) close() {
	close(em.chErr)
	atomic.StoreUint64(&em.isClosed, channelClosed)
}
