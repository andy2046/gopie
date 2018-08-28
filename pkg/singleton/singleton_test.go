package singleton

import (
	"testing"
)

func TestSingleton(t *testing.T) {
	s := New()
	s.Values["this"] = "that"

	s2 := New()
	if s2.Values["this"] != "that" {
		t.Fatal("wrong singleton implementation")
	}
}
