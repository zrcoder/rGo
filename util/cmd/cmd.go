package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"time"
)

const (
	OsWin    = "windows"
	OsLinux  = "linux"
	OsDarwin = "darwin"
)

// Executes the command and returns the result.
func Run(command string) (string, string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case OsWin:
		cmd = exec.Command(command)
	default:
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Run a command, kill it if timeout
func RunWithTimeout(timeout time.Duration, command string) (string, string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case OsWin:
		cmd = exec.Command(command)
	default:
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}

	done := make(chan error, 1)
	t := time.After(timeout)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case err = <-done:
		return stdout.String(), stderr.String(), err
	case <-t:
		cmd.Process.Kill()
		err = errors.New("command timed out")
		return stdout.String(), stderr.String(), err
	}
}
