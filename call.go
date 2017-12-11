package muts

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Printf is used to log and is default set to log.Printf
var PrintfFunc = log.Printf

var execCommand = exec.Command

// Call runs an operating system command composed of the parameters given.
// If only one parameter is given then interpret that as a single command line.
// Block until the command has finished.
// Calls Abort if the call failed.
func Call(params ...interface{}) int {
	r := Exec(NewExecOptions(params...))
	if !r.Ok() {
		Abort(r.Error)
	}
	return r.PID
}

// CallReturn runs an operating system command composed of the parameters given.
// If only one parameter is given then interpret that as a single command line.
// Block until the command has finished.
// Returns what is written to stdout or to stderror if an error was detected.
// The last newline(s) (if any) are stripped.
func CallReturn(params ...interface{}) (string, error) {
	errBuffer := new(bytes.Buffer)
	outBuffer := new(bytes.Buffer)
	opts := NewExecOptions(params...)
	opts.Stderr(errBuffer)
	opts.Stdout(outBuffer)
	r := Exec(opts)
	if !r.Ok() {
		return outBuffer.String() + "\n" + errBuffer.String(), errors.New(r.Error)
	}
	return strings.TrimRight(outBuffer.String(), "\n"), nil
}

// CallBackground runs an operating system command composed of the parameters given.
// If only one parameter is given then interpret that as a single command line.
// The output (stderr and stdout) of the result will be empty.
// Does not block ; this function returns immediately. Use the PID value of the result to handle the process.
func CallBackground(params ...interface{}) ExecResult {
	return Exec(NewExecOptions(params...).Wait(false))
}

// ExecResult holds the result of a Call.
type ExecResult struct {
	PID         int
	CommandLine string
	Error       string
	Stderr      string
	Stdout      string
}

// Ok returns whether the call was succesful. Only valid is the call was not run in the background.
func (r ExecResult) Ok() bool {
	return r.PID != 0
}

// ExecOptions is a parameter object for the Exec call
type ExecOptions struct {
	parameters []interface{}
	wait       bool
	force      bool
	output     io.Writer
	input      io.Reader
	errput     io.Writer
}

// Wait sets whether the call should wait for the command to complete. Default is true.
func (o *ExecOptions) Wait(w bool) *ExecOptions {
	o.wait = w
	return o
}

// Force determines if the program should continue when the cmd fails. If false the application will abort and defers
// are fired. Default is false.
func (o *ExecOptions) Force(f bool) *ExecOptions {
	o.force = f
	return o
}

// When Silent is set to true stdout and stderr will be discarded. Otherwise it is streamed as usual. Default is false.
// especially useful combined with force when you know you want re-runnable commands.
func (o *ExecOptions) Silent(s bool) *ExecOptions {
	if s {
		o.output = ioutil.Discard
		o.errput = ioutil.Discard
	}
	return o
}

// Stdout sets the writer for capturing the output produced by the command.
func (o *ExecOptions) Stdout(w io.Writer) *ExecOptions {
	o.output = w
	return o
}

// Stderr sets the writer for capturing the output produced by the command.
func (o *ExecOptions) Stderr(w io.Writer) *ExecOptions {
	o.errput = w
	return o
}

// Stderr sets the reader for accepting the input needed by the command.
func (o *ExecOptions) Stdin(r io.Reader) *ExecOptions {
	o.input = r
	return o
}

// Parameters sets the command and arguments. Can be a combination of values that are Stringers.
func (o *ExecOptions) Parameters(params ...interface{}) *ExecOptions {
	o.parameters = params
	return o
}

// NewExecOptions returns a new ExecOptions to be used in a Exec function call.
func NewExecOptions(params ...interface{}) *ExecOptions {
	return &ExecOptions{
		parameters: params,
		wait:       true,
		force:      false,
		output:     os.Stdout,
		input:      os.Stdin,
		errput:     os.Stderr,
	}
}

// Exec runs a shell command with parameters and settings from CallOptions
// Returns the process ID if not waiting and then 0 means there is a problem.
func Exec(options *ExecOptions) ExecResult {
	args := make([]string, len(options.parameters))
	for i, each := range options.parameters {
		args[i] = paramAsString(each)
	}
	if len(args) == 1 { // tokenize
		args = strings.Split(args[0], " ")
	}
	cmdline := strings.Join(args, " ")
	PrintfFunc("[sh -c] %s", cmdline)
	cmd := execCommand("sh", "-c", cmdline)
	cmd.Stdin = options.input

	cmd.Stdout = options.output
	cmd.Stderr = options.errput
	if options.wait {
		if err := cmd.Run(); err != nil && !options.force {
			return ExecResult{
				Error: err.Error(),
			}
		}
	} else {
		if err := cmd.Start(); err != nil {
			return ExecResult{
				Error: err.Error(),
			}
		}
	}
	if cmd.Process == nil {
		// if we don't know why
		return ExecResult{}
	}
	return ExecResult{PID: cmd.Process.Pid, CommandLine: cmdline}
}

func paramAsString(p interface{}) string {
	if s, ok := p.(string); ok {
		return s
	}
	if s, ok := p.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", p)
}
