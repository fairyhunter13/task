package task

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorManager_Success(t *testing.T) {
	m := NewErrorManager(0)

	var test string
	m.Run(func() error {
		test = "hello"
		return nil
	})
	var test2 string
	m.Run(func() error {
		test2 = "hi"
		return nil
	})

	err := m.Error()
	assert.EqualValues(t, test, "hello")
	assert.EqualValues(t, test2, "hi")
	assert.Nil(t, err)

	var k int
	m.Run(func() error {
		k = 2
		return nil
	})
	var k2 int64
	m.Run(func() error {
		k2 = 6
		return nil
	})

	err = m.Error()

	assert.EqualValues(t, k, 2)
	assert.EqualValues(t, k2, 6)
	assert.Nil(t, err)
}

func TestManager_EdgeCases(t *testing.T) {
	t.Run("No Init", func(*testing.T) {
		m := new(ErrorManager)
		m.Error()
	})
	t.Run("Single Error Job", func(t *testing.T) {
		m := NewErrorManager(0)

		var test string
		m.Run(func() error {
			test = "hello"
			return errors.New("there is an error here")
		})

		err := m.Error()
		assert.EqualValues(t, test, "hello")
		assert.NotNil(t, err)
	})
	t.Run("Sending On Closed Channel", func(t *testing.T) {
		m := NewErrorManager(0)

		var test string
		m.Run(func() error {
			test = "hello"
			return errors.New("there is an error here")
		})
		var test2 string
		m.Run(func() error {
			test2 = "hi"
			return errors.New("there is an error here")
		})
		m.Run(nil)

		err := m.Error()
		assert.EqualValues(t, test, "hello")
		assert.EqualValues(t, test2, "hi")
		assert.NotNil(t, err)
	})
	t.Run("Call Error Early", func(t *testing.T) {
		m := NewErrorManager(0)
		m.Error()
		m.ErrChan()

		var test string
		m.Run(func() error {
			test = "hello"
			return errors.New("there is an error here")
		})
		var test2 string
		m.Run(func() error {
			test2 = "hi"
			return nil
		})
		var test3 string
		m.Run(func() error {
			test3 = "hola"
			return nil
		})
		var test4 string
		m.Run(func() error {
			test4 = "ohayou"
			return nil
		})
		err := m.Error()
		assert.EqualValues(t, test, "hello")
		assert.EqualValues(t, test2, "hi")
		assert.EqualValues(t, test3, "hola")
		assert.EqualValues(t, test4, "ohayou")
		assert.NotNil(t, err)
	})
	t.Run("Using manual error channel", func(t *testing.T) {
		m := NewErrorManager(0)
		chanErr := m.ErrChan()

		var test string
		m.Run(func() error {
			test = "hello"
			return nil
		})
		var test2 string
		m.Run(func() error {
			test2 = "hi"
			return nil
		})
		var test3 string
		m.Run(func() error {
			test3 = "hola"
			return nil
		})
		var test4 string
		m.Run(func() error {
			test4 = "ohayou"
			return nil
		})

		m.WaitClose()
		for errTemp := range chanErr {
			assert.Nil(t, errTemp)
		}

		assert.EqualValues(t, test, "hello")
		assert.EqualValues(t, test2, "hi")
		assert.EqualValues(t, test3, "hola")
		assert.EqualValues(t, test4, "ohayou")
	})
}
