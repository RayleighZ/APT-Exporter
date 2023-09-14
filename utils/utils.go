package utils

import (
	"fmt"
	"os/exec"
)

func ExecAndGetRes(cmd string) string {
	command := exec.Command("bash", "-c", cmd)
	stdout, err := command.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(stdout)
}
