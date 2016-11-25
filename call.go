package muts

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

// Call runs a Command composed of the parameters given.
// If only one parameter is given then interpret that as a single command line.
// Abort the program if a call fails.
func Call(params ...interface{}) int { return WaitCall(true, false, params...) }

func CallForce(params ...interface{}) int { return WaitCall(true, true, params...) }

// WaitCall is the same as Call but has the option to wait for completion.
// Returns the process ID if not waiting and then 0 means there is a problem.
// The force parameter controls whether the WaitCall is aborted on error.
func WaitCall(wait bool, force bool, params ...interface{}) int {
	args := make([]string, len(params))
	for i, each := range params {
		args[i] = paramAsString(each)
	}
	if len(args) == 1 { // tokenize
		args = strings.Split(args[0], " ")
	}
	log.Println("sh -c", strings.Join(args, " "))
	cmd := execCommand("sh", "-c", strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if wait {
		if err := cmd.Run(); err != nil && !force {
			Fatalln(fmt.Sprintf("[run failed] %v in %s\n", err, Workspace))
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
