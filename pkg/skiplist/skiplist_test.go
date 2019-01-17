package skiplist

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"unsafe"
)

var testList *SkipList

func init() {
	testList = New()

	for i := 0; i <= 10000000; i++ {
		testList.Set(strconv.Itoa(i), int64(i))
	}

	var sl SkipList
	var el Element
	fmt.Printf("Sizeof(SkipList) = %v bytes Sizeof(Element) = %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
	fmt.Printf("Alignof(SkipList) = %v bytes Alignof(Element) = %v bytes\n", unsafe.Alignof(&sl), unsafe.Alignof(el))
}

func TestGetSet(t *testing.T) {
	list := New()
	n := 1000000
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := 0; i < n; i++ {
			list.Set(strconv.Itoa(i), int64(i))
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < n; i++ {
			list.Get(strconv.Itoa(i))
		}
		wg.Done()
	}()

	wg.Wait()
	if list.Len() != n {
		t.Fail()
	}
}

func TestCRUD(t *testing.T) {
	list := New()
	list.Set("A", 1)
	list.Set("B", 2)
	list.Set("C", 3)
	list.Set("D", 4)
	list.Set("E", 5)
	list.Set("C", 9)
	list.Remove("X")
	list.Remove("D")

	a := list.Get("A")
	b := list.Get("B")
	c := list.Get("C")
	d := list.Get("D")
	e := list.Get("E")
	x := list.Get("X")

	if a == nil || a.Key() != "A" || a.Value() != 1 {
		t.Fatal("wrong value", a)
	}
	if b == nil || b.Key() != "B" || b.Value() != 2 {
		t.Fatal("wrong value", b)
	}
	if c == nil || c.Key() != "C" || c.Value() != 9 {
		t.Fatal("wrong value", c)
	}
	if d != nil {
		t.Fatal("wrong value", d)
	}
	if e == nil || e.Key() != "E" || e.Value() != 5 {
		t.Fatal("wrong value", e)
	}
	if x != nil {
		t.Fatal("wrong value", x)
	}

}

func BenchmarkSet(b *testing.B) {
	b.ReportAllocs()
	list := New()

	for i := 0; i < b.N; i++ {
		list.Set(strconv.Itoa(i), int64(i))
	}
}

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		e := testList.Get(strconv.Itoa(i))
		if e == nil {
			b.Fatal("fail to Get")
		}
	}
}
