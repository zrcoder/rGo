package cmd

import (
	"os/exec"
	"time"
	"errors"
	"bytes"
)

// Executes the command and returns the result.
func Run(cmd string) (string, error) {
	command := exec.Command("/bin/sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Run a command, kill it if timeout
func RunWithTimeout(timeout time.Duration, command string, args ...string) (err error) {
	cmd := exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return
	}

	done := make(chan error, 1)
	t := time.After(timeout)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case err = <-done:
		return
	case <-t:
		cmd.Process.Kill()
		err = errors.New("command timed out")
		return
	}
	return
}
