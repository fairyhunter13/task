package task

import (
	"sync"

	"github.com/panjf2000/ants"
)

// Manager contains all the required tools to manage the task.
type Manager struct {
	wg *sync.WaitGroup
}

// NewManager initialize the task manager.
func NewManager() *Manager {
	m := new(Manager)
	m.init()
	return m
}

func (m *Manager) init() {
	if m.wg == nil {
		m.wg = new(sync.WaitGroup)
	}
}

// ClosureAnonym defines the anonymous function for the Run argument.
type ClosureAnonym func()

// Run runs the task in a separate go function.
func (m *Manager) Run(fn ClosureAnonym, opts ...OptionFunc) {
	if fn == nil {
		return
	}

	m.init()
	opt := new(Option).Assign(opts...)
	m.wg.Add(1)
	_ = ants.Submit(func() {
		defer m.wg.Done()

		if opt.UsePanicHandler {
			m.recoverPanic(fn)
		} else {
			fn()
		}
	})
}

func (m *Manager) recoverPanic(fn ClosureAnonym) {
	defer func() {
		_ = recover()
	}()
	fn()
}

// Wait blocks the current thread until the wg counter is zero.
func (m *Manager) Wait() {
	m.init()
	m.wg.Wait()
}
