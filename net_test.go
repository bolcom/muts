package muts

import "testing"

func TestPortLocal(t *testing.T) {
	PortRegistry = map[string]int{}
	*LocalUse = true
	i := Port("test", 8888)
	if got, want := i, 8888; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := PortRegistry["test"], 8888; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestPort(t *testing.T) {
	PortRegistry = map[string]int{}
	*LocalUse = false
	i := Port("test", 8888)
	if got, want := i, 8888; got == want {
		t.Errorf("got %v do not want %v", got, want)
	}
	if got, want := PortRegistry["test"], 8888; got == want {
		t.Errorf("got %v do not want %v", got, want)
	}
}
