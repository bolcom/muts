package muts

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// Abort is the function that is called if any error was detected.
// You can inject your own here but make sure to set it before calling any Defer(task).
var Abort = log.Fatalln

var tasks = map[string]func(){}

// LocalUse holds the value for the -local flag (default is false)
var LocalUse = flag.Bool("local", false, "Run all on your local machine")

// Task registers a function that can be called using a name as argument of the program (or via RunTask).
func Task(name string, f func()) {
	// for now, use a map of string->func, this may change
	tasks[name] = f
}

// RunTasks runs all the named tasks in order.
// If only one string is given then interpret that as a composition of names separated by spaces.
// After each task run, the current working directory is reset.
// Do not abort if a task is unknown.
func RunTasks(names ...string) {
	defer Chdir(Workspace)
	for _, each := range names {
		for _, name := range strings.Split(each, " ") {
			task, ok := tasks[name]
			if ok {
				PrintfFunc("\n----------------------\n task %q in %s\n----------------------\n", name, Workspace)
				Chdir(Workspace)
				task()
			} else {
				PrintfFunc("[RunTasks failed] unknown task %q", name)
			}
		}
	}
}

// RunTasksFromArgs reads the command line and runs each named task.
// A task name is an argument that does not start with a dash "-".
// Flags can be set before or after the task names.
// Boolean flags must have the form `-v=true`.
func RunTasksFromArgs() {
	defer deferList.run()
	if len(os.Args) == 1 {
		PrintTasks()
		fmt.Println()
		flag.PrintDefaults()
		return
	}
	for i := 1; i < len(os.Args); i++ {
		each := os.Args[i]
		if strings.HasPrefix(each, "-") {
			if !strings.Contains(each, "=") {
				i++ // next can be flag argument or a bool flag
			}
		} else {
			RunTasks(each)
		}
	}
}

// PrintTasks lists (in sort order) the names of all registered tasks.
func PrintTasks() {
	names := []string{}
	for each, _ := range tasks {
		names = append(names, each)
	}
	sort.Strings(names)
	for _, each := range names {
		fmt.Printf("\t%s\n", each)
	}
}
