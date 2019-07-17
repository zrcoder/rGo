package cmd

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tests := []struct {
		cmd string
	}{
		{"ls"},
		{"dir"},
		{"nslookup"},
		{"hhh"},
	}
	for _, test := range tests {
		t.Log(test.cmd)
		stdout, stderr, err := Run(test.cmd)
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
}

func TestRunWithTimeout(t *testing.T) {
	tests := []struct {
		cmd string
	}{
		{"nslookup"},
		{"hhh"},
		{"ipconfig"},
		{"ifconfig"},
	}

	for _, test := range tests {
		inner(test.cmd, 2*time.Nanosecond, t)
		inner(test.cmd, 20*time.Second, t)
	}
}

func inner(cmd string, timeout time.Duration, t *testing.T) {
	t.Log("command:", cmd, "; timeout:", timeout)
	stdout, stderr, err := RunWithTimeout(timeout, cmd)
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
