package muts

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Call runs a Command composed of the parameters given.
// If only one parameter is given then interpret that as a single command line.
// Abort the program if a call fails.
func Call(params ...string) int { return WaitCall(true, params...) }

// WaitCall is the same as Call but has the option to wait for completion.
// Returns the process ID if not waiting and then 0 means there is a problem.
func WaitCall(wait bool, params ...string) int {
	args := params
	if len(params) == 1 { // tokenize
		args = strings.Split(params[0], " ")
	}
	log.Println(strings.Join(args, " "))
	cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if wait {
		if err := cmd.Run(); err != nil {
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
