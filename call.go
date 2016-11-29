package muts

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

// Call runs a Command composed of the parameters given.
// If only one parameter is given then interpret that as a single command line.
// Abort the program if a call fails.
func Call(params ...interface{}) int {
	return Exec(NewExecOptions(params...))
}

// CallBackground runs a Command composed of the parameters given.
// If only one parameter is given then interpret that as a single command line.
// Return the PID of the process that was started or 0 if problem occurred
func CallBackground(params ...interface{}) int {
	return Exec(NewExecOptions(params...).Wait(false))
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

// Wait sets whether the call should proceed if the call fails. Default is false.
func (o *ExecOptions) Force(f bool) *ExecOptions {
	o.force = f
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
func Exec(options *ExecOptions) int {
	args := make([]string, len(options.parameters))
	for i, each := range options.parameters {
		args[i] = paramAsString(each)
	}
	if len(args) == 1 { // tokenize
		args = strings.Split(args[0], " ")
	}
	cmdline := strings.Join(args, " ")
	log.Println("sh -c", cmdline)
	cmd := execCommand("sh", "-c", cmdline)
	cmd.Stdin = options.input
	cmd.Stdout = options.output
	cmd.Stderr = options.errput
	if options.wait {
		if err := cmd.Run(); err != nil && !options.force {
			Fatalln(fmt.Sprintf("[run failed] %v -> %v, %v in %s\n", options.parameters, cmd.Args, err, Workspace))
		}
	} else {
		if err := cmd.Start(); err != nil {
			log.Println("[run on background failed] " + err.Error())
		}
	}
	if cmd.Process == nil {
		// if we don't know
		return 0
	}
	return cmd.Process.Pid
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
