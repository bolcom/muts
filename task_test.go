package muts

import (
	"os"
	"testing"
)

func TestRunTasksFromArgs(t *testing.T) {
	testHasRun := false
	Task("test", func() {
		testHasRun = true
	})
	os.Args = []string{"muts", "test"}
	RunTasksFromArgs()
	if !testHasRun {
		t.Fail()
	}
}
func TestRunTasksFromArgsWithFlag(t *testing.T) {
	testHasRun := false
	Task("test", func() {
		testHasRun = true
	})
	os.Args = []string{"muts", "-flag", "flag", "test"}
	RunTasksFromArgs()
	if !testHasRun {
		t.Fail()
	}
}
func TestRunTasksFromArgsWithFlagAssign(t *testing.T) {
	testHasRun := false
	Task("test", func() {
		testHasRun = true
	})
	os.Args = []string{"muts", "-flag=flag", "test"}
	RunTasksFromArgs()
	if !testHasRun {
		t.Fail()
	}
}

func TestRunTasksFromArgsWithFlagAssignWithEquals(t *testing.T) {
	testHasRun := false
	Task("test", func() {
		testHasRun = true
	})
	os.Args = []string{"muts", "-flag=flag=1", "test"}
	RunTasksFromArgs()
	if !testHasRun {
		t.Fail()
	}
}

func TestRunTasksFromArgsWithFlagAtEnd(t *testing.T) {
	testHasRun := false
	Task("test", func() {
		testHasRun = true
	})
	os.Args = []string{"muts", "test", "-flag", "flag"}
	RunTasksFromArgs()
	if !testHasRun {
		t.Fail()
	}
}
