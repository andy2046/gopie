package bitflag_test

import (
	"testing"

	. "github.com/andy2046/gopie/pkg/bitflag"
)

const (
	FLAG_A Flag = 1 << Flag(iota)
	FLAG_B
	FLAG_C
	FLAG_D
)

func TestExample(t *testing.T) {

	var flag Flag

	flag.SetAll(FLAG_A)
	flag.SetAll(FLAG_B, FLAG_C)
	flag.SetAll(FLAG_C | FLAG_D)

	flag.Reset()

	flag.SetAll(FLAG_A, FLAG_B, FLAG_C, FLAG_D)

	flag.ClearAll(FLAG_A)

	flag.ClearAll(FLAG_B, FLAG_C)

	if flag.AreAllSet(FLAG_A) {
		t.Fatal("A")
	}

	if flag.AreAllSet(FLAG_B) {
		t.Fatal("B")
	}

	if flag.AreAllSet(FLAG_C) {
		t.Fatal("C")
	}

	if !flag.AreAllSet(FLAG_D) {
		t.Fatal("D")
	}
}

func TestSet(t *testing.T) {

	var flag Flag

	flag.SetAll(FLAG_A)
	flag.SetAll(FLAG_B)

	if !flag.AreAllSet(FLAG_A) {
		t.Fail()
	}

	if !flag.AreAllSet(FLAG_B) {
		t.Fail()
	}

	flag.SetAll(FLAG_A, FLAG_B)

	if !flag.AreAllSet(FLAG_A, FLAG_B) {
		t.Fail()
	}

	if !flag.IsAnySet(FLAG_B) {
		t.Fail()
	}
}

func TestClear(t *testing.T) {

	var flag Flag

	flag.SetAll(FLAG_A, FLAG_B, FLAG_C)

	if !flag.AreAllSet(FLAG_A, FLAG_B, FLAG_C) {
		t.Fail()
	}

	flag.ClearAll(FLAG_B)

	if flag.AreAllSet(FLAG_B) {
		t.Fail()
	}

	if !flag.AreAllSet(FLAG_A, FLAG_C) {
		t.Fail()
	}
}

func TestToggle(t *testing.T) {

	var flag Flag

	flag.SetAll(FLAG_A, FLAG_B, FLAG_C)

	if !flag.AreAllSet(FLAG_A, FLAG_B, FLAG_C) {
		t.Fail()
	}

	flag.ToggleAll(FLAG_B)

	if flag.AreAllSet(FLAG_B) {
		t.Fail()
	}

	if !flag.AreAllSet(FLAG_A, FLAG_C) {
		t.Fail()
	}
}

func TestN(t *testing.T) {

	var flag Flag

	flag.Set(1)

	if !flag.IsSet(1) {
		t.Fail()
	}

	flag.Toggle(1)

	if flag.IsSet(1) {
		t.Fail()
	}

	flag.Toggle(1)
	flag.Clear(1)

	if flag.IsSet(1) {
		t.Fail()
	}
}
