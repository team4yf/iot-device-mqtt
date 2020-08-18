//Package kvstore the leveldb
package kvstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	Init("data.test")
	data, err := Get("config")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "foo", string(data))
}
