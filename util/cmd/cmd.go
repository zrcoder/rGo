package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"time"
)

// Executes the command and returns the result.
func Run(command string, arg ...string) (string, string, error) {
	cmd := exec.Command(command, arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Run a command, kill it if timeout
func RunWithTimeout(timeout time.Duration, command string, arg ...string) (string, string, error) {
	cmd := exec.Command(command, arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return "", "", err
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
		return "", "", err
	}
}

func BashRun(command string) (string, string, error) {
	return Run("/bin/bash", "-c", command)
}

func BashRunWithTimeout(timeout time.Duration, command string) (string, string, error) {
	return RunWithTimeout(timeout, "/bin/bash", "-c", command)
}
