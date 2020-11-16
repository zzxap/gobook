package main

import "testing"

func TestS1(t *testing.T) {
	if s1("abcdefgh") != 9 {
		t.Error(`s1("abcdefgh") != 9`)
	}

	if s1("") != 0 {
		t.Error(`s1("") != 0`)
	}
}
