package cmd

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tests := []struct {
		cmd string
	}{
		{"dir"},
		{"hhh"},
	}
	for _, test := range tests {
		t.Log(test.cmd)
		stdout, stderr, err := Run(test.cmd)
		printResult(t, stdout, stderr, err)
	}

	t.Log("ls -lah")
	stdout, stderr, err := Run("ls", "-lah")
	printResult(t, stdout, stderr, err)
}

func TestRunWithTimeout(t *testing.T) {
	tests := []struct {
		cmd string
	}{
		{"nslookup"},
		{"hhh"},
	}

	for _, test := range tests {
		inner(test.cmd, 2*time.Nanosecond, t)
		inner(test.cmd, 20*time.Second, t)
	}
}

func inner(cmd string, timeout time.Duration, t *testing.T) {
	t.Log("command:", cmd, "; timeout:", timeout)
	stdout, stderr, err := RunWithTimeout(timeout, cmd)
	printResult(t, stdout, stderr, err)
}

func TestBashRun(t *testing.T) {
	t.Log("ls -lah")
	stdout, stderr, err := BashRun("ls -lah")
	printResult(t, stdout, stderr, err)
}

func TestBashRunWithTimeout(t *testing.T) {
	t.Log("ls -lah; in 1 milli second")
	stdout, stderr, err := BashRunWithTimeout(time.Millisecond, "ls -lah")
	printResult(t, stdout, stderr, err)
}

func printResult(t *testing.T, stdout, stderr string, err error) {
	if stdout != "" {
		t.Logf("stdout:\n%s\n", stdout)
	}
	if stderr != "" {
		t.Logf("stderr:\n%s\n", stderr)
	}
	if err != nil {
		t.Log("err:", err)
	}
}
