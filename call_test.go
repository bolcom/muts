package muts

import (
	"os/exec"
	"testing"
)

func TestWaitCall(t *testing.T) {
	var args []string
	var cmd string
	execCommand = func(prog string, params ...string) *exec.Cmd {
		args = params
		cmd = prog
		return new(exec.Cmd)
	}
	defer func() { execCommand = exec.Command }()
	Exec(NewExecOptions("ls").Wait(true).Force(true))

	t.Log(cmd, args)
	if got, want := cmd, "sh"; got != want {
		t.Errorf("got %q want %q", got, want)
	}
	if got, want := args[0], "-c"; got != want {
		t.Errorf("got %q want %q", got, want)
	}
	if got, want := args[1], "ls"; got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCaptureOutput(t *testing.T) {
	out, err := CallReturn("date")
	if err != nil {
		t.Error("date output expected", err)
	}
	t.Log("stdout", out)
}
