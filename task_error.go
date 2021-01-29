package task

import (
	"sync/atomic"
)

const (
	channelClosed uint64 = 1
	channelOpen   uint64 = 0
)

// ErrorManager contains Manager but especially for error function.
type ErrorManager struct {
	chErr    chan error
	isClosed uint64
	counter  uint64
}

// NewErrorManager initialize the new error manager.
func NewErrorManager() *ErrorManager {
	em := new(ErrorManager)
	em.init()
	return em
}

func (em *ErrorManager) init() {
	if em.chErr == nil || em.isChannelClosed() {
		em.chErr = make(chan error)
		atomic.StoreUint64(&em.isClosed, channelOpen)
	}
}

func (em *ErrorManager) isChannelClosed() bool {
	return atomic.LoadUint64(&em.isClosed) == channelClosed
}

// ClosureErr defines closure that returns error.
type ClosureErr func() (err error)

// Run runs the closure error function.
func (em *ErrorManager) Run(fn ClosureErr) {
	em.init()
	atomic.AddUint64(&em.counter, 1)
	go func() {
		em.chErr <- fn()
	}()
}

// ErrChan returns the receiving error channel of this error manager.
func (em *ErrorManager) ErrChan() <-chan error {
	em.init()
	return em.chErr
}

// Error returns the first error from the Run execution of the fn closure.
func (em *ErrorManager) Error() (err error) {
	em.init()
	defer em.close()

	if em.isJobDone() {
		return
	}

	for {
		errTemp, more := <-em.chErr
		atomic.AddUint64(&em.counter, ^uint64(0))
		if !more {
			break
		}

		if err == nil {
			err = errTemp
		}

		if em.isJobDone() {
			break
		}
	}
	return
}

func (em *ErrorManager) isJobDone() bool {
	return atomic.LoadUint64(&em.counter) == 0
}

func (em *ErrorManager) close() {
	close(em.chErr)
	atomic.StoreUint64(&em.isClosed, channelClosed)
}
