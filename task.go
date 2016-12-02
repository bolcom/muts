package muts

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// Tasks holds a mapping between task names and the function to call for running it.
var Tasks = map[string]func(){}

// LocalUse holds the value for the -local flag (default is false)
var LocalUse = flag.Bool("local", false, "Run all on your local machine")

// RunTasks runs all the named tasks in order.
// If only one string is given then interpret that as a composition of names separated by spaces.
// After each task run, the current working directory is reset.
// Do not abort if a task is unknown.
func RunTasks(names ...string) {
	defer Chdir(Workspace)
	for _, each := range names {
		for _, name := range strings.Split(each, " ") {
			task, ok := Tasks[name]
			if ok {
				log.Printf("task [%s] in %s\n", name, Workspace)
				Chdir(Workspace)
				task()
			} else {
				log.Printf("[RunTasks failed] unknown task %q", name)
			}
		}
	}
}

// RunTasksFromArgs reads the command line and runs each named task.
// A task name is an argument that does not start with a dash "-"
func RunTasksFromArgs() {
	if len(os.Args) == 1 {
		PrintTasks()
		fmt.Println()
		flag.PrintDefaults()
		return
	}
	for i := 1; i < len(os.Args); i++ {
		each := os.Args[i]
		if !strings.HasPrefix(each, "-") {
			RunTasks(os.Args[i])
		}
	}
}

// PrintTasks lists (in sort order) the names of all registered tasks.
func PrintTasks() {
	names := []string{}
	for each, _ := range Tasks {
		names = append(names, each)
	}
	sort.Strings(names)
	for _, each := range names {
		fmt.Printf("\t%s\n", each)
	}
}
