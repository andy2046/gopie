package lru

import (
	"fmt"
	"testing"
)

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

var getTests = []struct {
	name       string
	keyToAdd   interface{}
	keyToGet   interface{}
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
	{"simple_struct_hit", simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
	{"simeple_struct_miss", simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
	{"complex_struct_hit", complexStruct{1, simpleStruct{2, "three"}},
		complexStruct{1, simpleStruct{2, "three"}}, true},
}

func TestGet(t *testing.T) {
	for _, tt := range getTests {
		lru := New(0)
		lru.Add(tt.keyToAdd, 1234)
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestRemove(t *testing.T) {
	lru := New(0)
	lru.Add("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Remove("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed entry")
	}
}

func TestPurge(t *testing.T) {
	purgedKeys := make([]interface{}, 0)
	onPurgedFun := func(key interface{}, value interface{}) {
		purgedKeys = append(purgedKeys, key)
	}

	lru := New(20)
	lru.OnPurged = onPurgedFun
	for i := 0; i < 22; i++ {
		lru.Add(fmt.Sprintf("myKey%d", i), 1234)
	}

	if len(purgedKeys) != 2 {
		t.Fatalf("got %d evicted keys; want 2", len(purgedKeys))
	}
	if purgedKeys[0] != interface{}("myKey0") {
		t.Fatalf("got %v in first evicted key; want %s", purgedKeys[0], "myKey0")
	}
	if purgedKeys[1] != interface{}("myKey1") {
		t.Fatalf("got %v in second evicted key; want %s", purgedKeys[1], "myKey1")
	}
}

func BenchmarkLRU(b *testing.B) {
	purgedKeys := make([]interface{}, 0)
	onPurgedFun := func(key interface{}, value interface{}) {
		purgedKeys = append(purgedKeys, key)
	}
	n, m := 20, 40
	lru := New(n)
	lru.OnPurged = onPurgedFun
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Add(fmt.Sprintf("myKey%d", i%m), 1234)
	}
}
