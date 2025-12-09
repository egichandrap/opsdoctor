package svcchecker

import (
	"os/exec"
)

func CheckService(name string) (string, error) {
	out, err := exec.Command("systemctl", "status", name).CombinedOutput()
	return string(out), err
}
