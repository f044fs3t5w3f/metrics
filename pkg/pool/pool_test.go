package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	state string
}

func (t *testStruct) Reset() {
	t.state = "reset"
}

func newTestStruct() *testStruct {
	return &testStruct{
		state: "new",
	}
}

func TestNew_wasReset(t *testing.T) {
	pool := New(newTestStruct)
	obj := pool.Get()
	t.Run("New object", func(t *testing.T) {
		assert.Equal(t, "new", obj.state)
	})

	pool.Put(obj)
	obj = pool.Get()

	t.Run("Object was reset", func(t *testing.T) {
		assert.Equal(t, "reset", obj.state)
	})
}
