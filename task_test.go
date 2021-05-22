package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_Success(t *testing.T) {
	m := NewManager()

	var test string
	m.Run(func() {
		test = "hello"
	})
	var test2 string
	m.Run(func() {
		test2 = "hi"
	})
	m.Run(nil)

	m.Wait()

	assert.EqualValues(t, test, "hello")
	assert.EqualValues(t, test2, "hi")

	var k int
	m.Run(func() {
		k = 2
	})
	var k2 int64
	m.Run(func() {
		k2 = 6
	})

	m.Wait()

	assert.EqualValues(t, k, 2)
	assert.EqualValues(t, k2, 6)
}

func TestManager_NoInit(t *testing.T) {
	m := new(Manager)

	var test string
	m.Run(func() {
		test = "hello"
	})
	var test2 string
	m.Run(func() {
		test2 = "hi"
	})

	m.Wait()

	assert.EqualValues(t, test, "hello")
	assert.EqualValues(t, test2, "hi")
}

func TestManager_Panic(t *testing.T) {
	m := NewManager()

	var test string
	m.Run(func() {
		test = "hello"
		panic("hello")
	}, WithPanicHandler(true))
	var test2 string
	m.Run(func() {
		test2 = "hi"
	})
	m.Run(nil)

	m.Wait()

	assert.EqualValues(t, test, "hello")
	assert.EqualValues(t, test2, "hi")
}
