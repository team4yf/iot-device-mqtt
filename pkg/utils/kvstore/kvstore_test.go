//Package kvstore the leveldb
package kvstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	Init("data.test")
	data, exists := Get("config")
	assert.False(t, exists, "exists should be false")
	assert.Equal(t, "", string(data))
}
func TestPut(t *testing.T) {
	Init("data.test")
	data, _ := Get("config")
	err := Put("config", []byte("foo3"))
	assert.Nil(t, err, "err should be nil")

	data, _ = Get("config")
	assert.Equal(t, "foo2", string(data))
}

func TestPutObject(t *testing.T) {
	Init("data.test")
	defer Close()
	err := PutObject("config", map[string]string{
		"foo": "bar",
	})
	assert.Nil(t, err, "err should be nil")
	obj := make(map[string]string)
	_ = GetObject("config", &obj)
	assert.Equal(t, "bcd", obj["foo"])
}
