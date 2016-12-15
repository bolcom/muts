package muts

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
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
	if !strings.Contains(out, fmt.Sprint(time.Now().Year())) {
		t.Errorf("got %q", out)
	}
}

func TestExecNonSilent(t *testing.T) {
	outBuffer := new(bytes.Buffer)
	result := Exec(NewExecOptions("echo TEST").Stdout(outBuffer).Silent(false))
	if !strings.Contains(outBuffer.String(), "TEST") {
		t.Errorf("got [%q] wanted at least [%q]", result, "TEST")
	}

	outBuffer = new(bytes.Buffer)
	result = Exec(NewExecOptions("echo TEST").Stdout(outBuffer).Silent(true))
	if strings.Contains(outBuffer.String(), "TEST") {
		t.Errorf("got [%q] wanted at least [%q]", result, "TEST")
	}

}
func TestExecSilentError(t *testing.T) {
	errBuffer := new(bytes.Buffer)
	result := Exec(NewExecOptions("echo TEST", " 1>&2 ").Stderr(errBuffer).Silent(false))
	if !strings.Contains(errBuffer.String(), "TEST") {
		t.Errorf("got [%q] wanted at least [%q]", result, "TEST")
	}

	errBuffer = new(bytes.Buffer)
	result = Exec(NewExecOptions("echo TEST", " 1>&2 ").Stderr(errBuffer).Silent(true))
	if strings.Contains(errBuffer.String(), "TEST") {
		t.Errorf("got [%q] wanted at least [%q]", result, "TEST")
	}
}
