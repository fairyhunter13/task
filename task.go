package task

import "sync"

// Manager contains all the required tools to manage the task.
type Manager struct {
	wg *sync.WaitGroup
}

// New initialize the task manager.
func New() *Manager {
	return &Manager{
		wg: new(sync.WaitGroup),
	}
}

func (m *Manager) initWg() {
	if m.wg == nil {
		m.wg = new(sync.WaitGroup)
	}
}

// AnonymousClosure defines the anonymous function for the Run argument.
type AnonymousClosure func()

// Run runs the task in a new go function.
func (m *Manager) Run(fn AnonymousClosure) {
	m.initWg()
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		fn()
	}()
}

// Wait blocks the current thread until the wg is zero.
func (m *Manager) Wait() {
	m.initWg()
	m.wg.Wait()
}
