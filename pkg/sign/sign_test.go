package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSignFunc(t *testing.T) {
	key := "cat"
	sign := GetSignFunc(key)
	t.Run("Simple Sign", func(t *testing.T) {
		hash := sign([]byte("Dog"))
		assert.Equal(t, []byte("\xe3\x10\u07ba\xca]S\x14Vǘ#\xfd\x18d^\xe4_\xf6NXF\xed\x01V\xe4\xe2\xd0\xec\xab\xce\xc4"), hash)
	})
}
