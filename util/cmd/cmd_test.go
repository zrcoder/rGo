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
		output, err := Run(test.cmd)
		if err != nil {
			t.Log(err)
		} else {
			t.Log(output)
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
		inner(test.cmd, 2, t)
		inner(test.cmd, 20*time.Second, t)
	}
}

func inner(cmd string, timeout time.Duration, t *testing.T) {
	t.Log("command:", cmd, "; timeout:", timeout)
	err := RunWithTimeout(timeout, cmd)
	if err != nil {
		t.Log("error hanppens:", err)
	} else {
		t.Log("not timeout!")
	}
}
