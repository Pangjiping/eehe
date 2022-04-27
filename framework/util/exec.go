package util

import (
	"os"
	"syscall"
)

// GetExecDir gets current project exec directory.
func GetExecDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir + "/"
}

// CheckProcessExist check if process PID exists.
func CheckProcessExist(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal(0) to this process, it exists if return nil
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
