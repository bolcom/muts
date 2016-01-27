package muts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspace(t *testing.T) {
	if len(Workspace) == 0 {
		t.Fail()
	}
}

func TestCreateFileWith(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "muts.txt")
	CreateFileWith(tmp, "garbage")
	if _, err := os.Stat(tmp); err != nil {
		t.Error("unable to create file", err)
	}
}

func TestCreateFileWith_Fail(t *testing.T) {
	caught := false
	Fatalln = func(args ...interface{}) {
		caught = true
	}
	CreateFileWith("/usr", "garbage")
	if !caught {
		t.Fatalf("expected fatal error")
	}
}
