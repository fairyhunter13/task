package task

import (
	"sync"
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

// Run runs the task in a new go function.
func (m *Manager) Run(fn ClosureAnonym) {
	m.init()
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		fn()
	}()
}

// Wait blocks the current thread until the wg counter is zero.
func (m *Manager) Wait() {
	m.init()
	m.wg.Wait()
}
