package ringhash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestAddNode(t *testing.T) {
	r := New()

	r.AddNode("127.0.0.1:80")
	if len(r.hashes) != r.replicas {
		t.Fatal("wrong vnodes number")
	}
}

func TestGetNode(t *testing.T) {
	r := New()

	r.AddNode("127.0.0.1:80")
	node, err := r.GetNode("127.0.0.1:80")
	if err != nil {
		t.Fatal(err)
	}

	if node != "127.0.0.1:80" {
		t.Fatalf("wrong node, expected 127.0.0.1:80, got %v\n", node)
	}
}

func TestRemoveNode(t *testing.T) {
	r := New()

	r.AddNode("127.0.0.1:80")
	r.RemoveNode("127.0.0.1:80")

	if len(r.hashes) != 0 && len(r.hashKeyMap) != 0 {
		t.Fatal(("remove not working"))
	}
}

func TestGetLeastNode(t *testing.T) {
	option := func(c *Config) error {
		c.BalancingFactor = 1.02
		return nil
	}
	r := New(option)

	r.AddNode("127.0.0.1:80")
	r.AddNode("192.168.0.1:80")
	r.AddNode("10.0.0.1:80")

	for i := 0; i < 100; i++ {
		node, err := r.GetLeastNode("192.168.0.1:81")
		if err != nil {
			t.Fatal(err)
		}
		r.Add(node)
	}

	for k, v := range r.Loads() {
		if v > r.MaxLoad() {
			t.Fatalf("node %s is overloaded, %d > %d\n", k, v, r.MaxLoad())
		}
	}
	fmt.Println("Max load per node ->", r.MaxLoad())
	fmt.Println(r.Loads())
}

func TestAddDone(t *testing.T) {
	r := New()

	r.AddNode("127.0.0.1:80")
	r.AddNode("192.168.0.1:80")

	node, err := r.GetLeastNode("192.168.0.1:81")
	if err != nil {
		t.Fatal(err)
	}

	r.Add(node)
	if r.keyLoadMap[node].Load != 1 {
		t.Fatalf("load for node %s should be 1\n", node)
	}

	r.Done(node)
	if r.keyLoadMap[node].Load != 0 {
		t.Fatalf("load for node %s should be 0\n", node)
	}
}

func BenchmarkGetNode(b *testing.B) {
	r := New()
	for i := 0; i < 10; i++ {
		r.AddNode("start" + strconv.Itoa(i))
	}
	tt := []struct {
		key string
	}{
		{"test"},
		{"test1"},
		{"test2"},
		{"test3"},
		{"test4"},
		{"test5"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := tt[i%len(tt)]
		r.GetNode(t.key)
	}
}

func BenchmarkGetLeastNode(b *testing.B) {
	r := New()
	for i := 0; i < 10; i++ {
		r.AddNode("start" + strconv.Itoa(i))
	}
	tt := []struct {
		key string
	}{
		{"test"},
		{"test1"},
		{"test2"},
		{"test3"},
		{"test4"},
		{"test5"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := tt[i%len(tt)]
		r.GetLeastNode(t.key)
	}
}

func BenchmarkAddRemoveNode(b *testing.B) {
	r := New()
	for i := 0; i < 10; i++ {
		r.AddNode("start" + strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.AddNode("foo" + strconv.Itoa(i))
		r.RemoveNode("foo" + strconv.Itoa(i))
	}
}

func BenchmarkAddDone(b *testing.B) {
	r := New()
	for i := 0; i < 10; i++ {
		r.AddNode("start" + strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node, _ := r.GetLeastNode("start" + strconv.Itoa(i))
		r.Add(node)
		r.Done(node)
	}
}
