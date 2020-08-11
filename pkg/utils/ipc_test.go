package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	out, err := RunCmd("ls -l /home/yf | grep yunplus | awk '{print $9}'")

	assert.Nil(t, err, "should not error")

	fmt.Printf("out: \n%v", (string)(out))
}
