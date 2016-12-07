package muts

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Abort is the function that is called if any error was detected. You can inject your own here
var Abort = func(v ...interface{}) {
	panic(fmt.Sprint(v...))
}

// Workspace holds the current working directory on startup.
var Workspace, _ = os.Getwd()

// CreateFileWith does what you think it should do.
func CreateFileWith(filename, contents string) {
	f, err := os.Create(filename)
	if err != nil {
		Abort("CreateFileWith failed:", err)
	}
	defer f.Close()
	_, err = io.WriteString(f, contents)
	if err != nil {
		Abort("CreateFileWith failed:", err)
	}
	log.Printf("written %d bytes to %s\n", len(contents), filename) // show absolute name
}

// Setenv wraps the os one to check and log it
func Setenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		Abort("Setenv failed:", err)
	}
	log.Println(key, "=", value)
}

// Chdir wraps the os one to check and log it
func Chdir(whereto string) {
	here, err := os.Getwd()
	if err != nil {
		Abort("Chdir failed:", err)
	}
	if here == whereto {
		return
	}
	abs, err := filepath.Abs(whereto)
	if err != nil {
		Abort("Chdir failed:", err)
	}
	if here == abs {
		return
	}
	err = os.Chdir(whereto)
	if err != nil {
		Abort("Chdir failed:", err)
	}
	log.Printf("changed workdir: [%s] -> [%s]", here, abs)
}
