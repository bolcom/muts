package muts

import (
	"log"
	"os"
	"path/filepath"
)

// Fatalln is the function that is called if any error was detected. You can inject your own here
var Fatalln = log.Fatalln

// Workspace holds the current working directory on startup.
var Workspace, _ = os.Getwd()

// CreateFileWith does what you think it should do.
func CreateFileWith(filename, contents string) {
	f, err := os.Create(filename)
	if err != nil {
		Fatalln("CreateFileWith failed:", err)
	}
	defer f.Close()
	n, err := f.WriteString(contents)
	if err != nil {
		Fatalln("CreateFileWith failed:", err)
	}
	log.Printf("written %d/%d bytes to %s\n", n, len(contents), filename) // show absolute name
}

// Setenv wraps the os one to check and log it
func Setenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		Fatalln("Setenv failed:", err)
	}
	log.Println(key, "=", value)
}

// Chdir wraps the os one to check and log it
func Chdir(whereto string) {
	here, err := os.Getwd()
	if err != nil {
		Fatalln("Chdir failed:", err)
	}
	if here == whereto {
		return
	}
	abs, err := filepath.Abs(whereto)
	if err != nil {
		Fatalln("Chdir failed:", err)
	}
	if here == abs {
		return
	}
	err = os.Chdir(whereto)
	if err != nil {
		Fatalln("Chdir failed:", err)
	}
	log.Printf("changed workdir: [%s] -> [%s]", here, abs)
}
