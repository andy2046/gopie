package dll

import (
	"fmt"
	"sync"
	"testing"
)

var (
	m    = make(map[int]bool)
	lock sync.Mutex
)

func TestDListEmpty(t *testing.T) {
	l := New()
	if !l.Empty() {
		t.Error("new list is empty")
	}

	var elm *Element

	elm = l.PushLeft(1)
	if l.Empty() {
		t.Error("list is not empty")
	}
	l.Remove(elm)
	if !l.Empty() {
		t.Error("list is empty")
	}

	elm = l.PushRight(2)
	if l.Empty() {
		t.Error("list is not empty")
	}
	l.Remove(elm)
	if !l.Empty() {
		t.Error("list is empty")
	}

	elm = l.PushLeft(3)
	if l.Empty() {
		t.Error("list is not empty")
	}
	l.PopRight()
	if !l.Empty() {
		t.Error("list is empty")
	}

	elm = l.PushRight(4)
	if l.Empty() {
		t.Error("list is not empty")
	}
	l.PopLeft()
	if !l.Empty() {
		t.Error("list is empty")
	}
}

func TestDListNext(t *testing.T) {
	l := New()
	start := make(chan struct{})
	n := 100

	var init *Element
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		<-start
		removeCnt := 0
		for i := 0; i < n; i++ {
			e := l.PushRight(i + n)
			fmt.Printf("pushRight %v\n", i+n)

			if i == 0 {
				init = e
				continue
			}

			if removeCnt < n/4 {
				r := l.Remove(e)
				removeCnt++
				if r != nil {
					fmt.Printf("remove %v\n", r)
				}
			}
		}
		wg.Done()
		fmt.Println("pushRight done")
	}()

	close(start)
	wg.Wait()

	nxt := l.Next(init)
	fmt.Printf("Next of %v is %v\n", init.Value, nxt.Value)

	for {
		v := l.PopLeft()
		if v == nil {
			fmt.Println()
			break
		}
		fmt.Printf("%v.", v)
	}

	if !l.Empty() {
		t.Error("list is empty")
	}
}

func TestDListPrev(t *testing.T) {
	l := New()
	start := make(chan struct{})
	n := 100

	var init *Element
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		<-start
		removeCnt := 0
		for i := 0; i < n; i++ {
			e := l.PushLeft(i + n)
			fmt.Printf("pushRight %v\n", i+n)

			if i == 0 {
				init = e
				continue
			}

			if removeCnt < n/4 {
				r := l.Remove(e)
				removeCnt++
				if r != nil {
					fmt.Printf("remove %v\n", r)
				}
			}
		}
		wg.Done()
		fmt.Println("pushRight done")
	}()

	close(start)
	wg.Wait()

	nxt := l.Prev(init)
	fmt.Printf("Prev of %v is %v\n", init.Value, nxt.Value)

	for {
		v := l.PopLeft()
		if v == nil {
			fmt.Println()
			break
		}
		fmt.Printf("%v.", v)
	}

	if !l.Empty() {
		t.Error("list is empty")
	}
}

func TestDListPL(t *testing.T) {
	l := New()
	start := make(chan struct{})
	n := 100

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		<-start
		removeCnt := 0
		for i := 0; i < n; i++ {
			e := l.PushRight(i + n)
			fmt.Printf("pushRight %v\n", i+n)
			if removeCnt < n/4 {
				r := l.Remove(e)
				removeCnt++
				if r != nil {
					// vv := r.(int)
					// record(t, vv)
					fmt.Printf("remove %v\n", r)
				}
			}
		}
		wg.Done()
		fmt.Println("pushRight done")
	}()

	go func() {
		<-start
		for i := 0; i < n/2; i++ {
			v := l.PopLeft()
			// if v != nil {
			// 	vv := v.(int)
			// 	record(t, vv)
			// }
			fmt.Printf("popLeft %v\n", v)
		}
		wg.Done()
		fmt.Println("popLeft done")
	}()

	close(start)
	wg.Wait()
	for {
		v := l.PopLeft()
		if v == nil {
			fmt.Println()
			break
		}
		fmt.Printf("%v.", v)
	}

	if !l.Empty() {
		t.Error("list is empty")
	}
}

func TestDListPR(t *testing.T) {
	l := New()
	start := make(chan struct{})
	n := 100

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		<-start
		removeCnt := 0
		for i := 0; i < n; i++ {
			e := l.PushLeft(i)
			fmt.Printf("pushLeft %v\n", i)
			if removeCnt < n/4 {
				r := l.Remove(e)
				removeCnt++
				if r != nil {
					// vv := r.(int)
					// record(t, vv)
					fmt.Printf("remove %v\n", r)
				}
			}
		}
		wg.Done()
		fmt.Println("pushLeft done")
	}()

	go func() {
		<-start
		for i := 0; i < n/2; i++ {
			v := l.PopRight()
			// if v != nil {
			// 	vv := v.(int)
			// 	record(t, vv)
			// }
			fmt.Printf("popRight %v\n", v)
		}
		wg.Done()
		fmt.Println("popRight done")
	}()

	close(start)
	wg.Wait()
	for {
		v := l.PopLeft()
		if v == nil {
			fmt.Println()
			break
		}
		fmt.Printf("%v.", v)
	}

	if !l.Empty() {
		t.Error("list is empty")
	}
}

func record(t *testing.T, v int) {
	lock.Lock()
	defer lock.Unlock()

	if _, existed := m[v]; existed {
		t.Fatalf("duplicated %v", v)
	}

	m[v] = true
}
