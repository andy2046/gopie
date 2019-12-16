package nocopy_test

import (
	. "github.com/andy2046/gopie/pkg/nocopy"
	"testing"
)

type MyStruct struct {
	noCopy NoCopy
}

func TestNoCopy(t *testing.T) {
	// go vet fails
	var m1 MyStruct
	m2 := m1
	var m3 = m1
	m2 = m1
	_, _ = m2, m3
	t.Log("go vet fails here")
}
