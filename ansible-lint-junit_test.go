package main

import "testing"

func TestVersion(t *testing.T) {
	want := AppVersion
	if got := Version(); got != want {
		t.Errorf("Version() = %q, want %q", got, want)
	}
}
