package main

import (
	"testing"
)

func TestParseArgsSpaces(t *testing.T) {
	out := parseArgs("account spaces")
	if len(out) != 2 {
		t.Fatalf("Output should have been length 2, was length %v", len(out))
	}
}

func TestParseArgsVariadic(t *testing.T) {
	out := parseArgs("account", "spaces")
	if len(out) != 2 {
		t.Fatalf("Output should have been length 2, was length %v", len(out))
	}
}
func TestParseArgsSomeSpacesSomeVariadic(t *testing.T) {
	out := parseArgs("account space", "mySpace")
	if len(out) != 3 {
		t.Fatalf("Output should have been length 3, was length %v", len(out))
	}
}
