package task

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorManager_Success(t *testing.T) {
	m := NewErrorManager(5)

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
	t.Run("Sending On Closed Channel", func(t *testing.T) {
		m := NewErrorManager(1)

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

		err := m.Error()
		assert.EqualValues(t, test, "hello")
		assert.EqualValues(t, test2, "")
		assert.NotNil(t, err)
	})
	t.Run("Buffer Size 0", func(t *testing.T) {

	})
}
