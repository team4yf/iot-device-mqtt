package utils

import (
	"os/exec"
)

func RunCmd(shellCmd string) ([]byte, error) {
	out, err := exec.Command("bash", "-c", shellCmd).Output()
	return out, err
}
