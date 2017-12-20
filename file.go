package muts

import (
	"log"
	"os"
	"path/filepath"
)

// Workspace holds the current working directory on startup.
var Workspace, _ = os.Getwd()

// CreateFileWith does what you think it should do.
func CreateFileWith(filename, contents string) {
	f, err := os.Create(filename)
	if err != nil {
		Abort("CreateFileWith failed:", err)
	}
	defer f.Close()
	n, err := f.WriteString(contents)
	if err != nil {
		Abort("CreateFileWith failed:", err)
	}
	log.Printf("written %d/%d bytes to %s\n", n, len(contents), filename) // show absolute name
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
	PrintfFunc("changed workdir: [%s] -> [%s]", here, abs)
}

// Mkdir wraps os.MkdirAll to check and log it.
func Mkdir(path string) {
	abs, err := filepath.Abs(path)
	err = os.MkdirAll(abs, os.ModePerm)
	if err != nil {
		Abort("Mkdir failed:", err)
	}
	PrintfFunc("created dir: [%s]", abs)
}
