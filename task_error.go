package task

import (
	"sync/atomic"
)

const (
	// DefaultBufferSize is the default lenght of the channel error manager
	DefaultBufferSize int    = 5
	channelClosed     uint64 = 1
	channelOpen       uint64 = 0
)

// ErrorManager contains Manager but especially for error function.
type ErrorManager struct {
	chErr      chan error
	bufferSize int
	isClosed   uint64
}

// NewErrorManager initialize the new error manager.
func NewErrorManager(bufferSize int) *ErrorManager {
	em := new(ErrorManager)
	em.init(bufferSize)
	return em
}

func (em *ErrorManager) init(bufferSize int) {
	if bufferSize == 0 {
		bufferSize = DefaultBufferSize
	}
	em.bufferSize = bufferSize

	if em.chErr == nil || em.isChannelClosed() {
		em.chErr = make(chan error, em.bufferSize)
		atomic.StoreUint64(&em.isClosed, channelOpen)
	}
}

func (em *ErrorManager) isChannelClosed() bool {
	return em.isClosed == channelClosed
}

// ClosureErr defines closure that returns error.
type ClosureErr func() (err error)

// Run runs the closure error function.
func (em *ErrorManager) Run(fn ClosureErr) {
	em.init(em.bufferSize)
	go func() {
		em.chErr <- fn()
	}()
}

// ErrChan returns the receiving error channel of this error manager.
func (em *ErrorManager) ErrChan() <-chan error {
	em.init(em.bufferSize)
	return em.chErr
}

// Error returns the first error from the Run execution of the fn closure.
func (em *ErrorManager) Error() (err error) {
	em.init(em.bufferSize)
	defer em.close()

	select {
	case err = <-em.chErr:
	default:
		return
	}

	for {
		var (
			errTemp error
			more    bool
		)
		errTemp, more = <-em.chErr
		if !more || em.isJobDone() {
			break
		}

		if err != nil {
			continue
		}

		err = errTemp
	}
	return
}

func (em *ErrorManager) isJobDone() bool {
	return len(em.chErr) == 0
}

func (em *ErrorManager) close() {
	close(em.chErr)
	atomic.StoreUint64(&em.isClosed, channelClosed)
}
