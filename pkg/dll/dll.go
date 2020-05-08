// Package dll provides a lock-free implementation of doubly linked list.
package dll

type (
	// List represents a doubly linked list.
	List struct {
		head *Element
		tail *Element
	}

	// Element is an element of a linked list.
	Element struct {
		// Next and previous pointers in the doubly-linked list of elements.
		next, prev *atomicMarkableReference

		// The value stored with this element.
		Value interface{}
	}
)

// New returns an initialized list.
func New() *List { return new(List).Init() }

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.head = &Element{
		prev: newAtomicMarkableReference(nil, false),
		next: newAtomicMarkableReference(nil, false),
	}
	l.tail = &Element{
		prev: newAtomicMarkableReference(nil, false),
		next: newAtomicMarkableReference(nil, false),
	}

	l.head.next.compareAndSet(nil, l.tail, false, false)
	l.tail.prev.compareAndSet(nil, l.head, false, false)

	return l
}

// Empty returns true if list l is empty, false otherwise
func (l *List) Empty() bool {
	t := l.head.next.getReference()
	h := l.tail.prev.getReference()

	if t == l.tail && h == l.head {
		return true
	}

	return false
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value if succeed, nil otherwise
func (l *List) Remove(e *Element) interface{} {
	node := e
	if node == nil {
		return nil
	}

	if node == l.tail || node == l.head {
		return nil
	}

	for {
		removed, nodeNext := node.next.get()
		if removed {
			return nil
		}
		if node.next.compareAndSet(nodeNext, nodeNext, false, true) {
			for {
				removed2, nodePrev := node.prev.get()
				if removed2 || node.prev.compareAndSet(nodePrev, nodePrev, false, true) {
					break
				}
			}

			l.correctPrev(node.prev.getReference(), nodeNext)
			return node.Value
		}
	}
}

func (l *List) correctPrev(prev, node *Element) *Element {
	var lastLink *Element

	for {
		removed, link1 := node.prev.get()
		if removed {
			break
		}
		removed2, prevNext := prev.next.get()
		if removed2 {
			if lastLink != nil {
				setMark(prev.prev)
				lastLink.next.compareAndSet(prev, prevNext, false, false)
				prev = lastLink
				lastLink = nil
				continue
			}
			prevNext = prev.prev.getReference()
			prev = prevNext
			continue
		}

		if prevNext != node {
			lastLink = prev
			prev = prevNext
			continue
		}

		if node.prev.compareAndSet(link1, prev, removed, false) {
			if prev.prev.isMarked() {
				continue
			}
			break
		}
	}

	return prev
}

func setMark(node *atomicMarkableReference) {
	for {
		removed, link := node.get()
		if removed || node.compareAndSet(link, link, false, true) {
			break
		}
	}
}

// PopLeft returns the first element of list l or nil if the list is empty.
func (l *List) PopLeft() interface{} {
	prev := l.head

	for {
		node := prev.next.getReference()
		// deque is empty
		if node == l.tail {
			return nil
		}

		removed, nodeNext := node.next.get()
		// concurrent pop started to delete this node, help it, then continue
		if removed {
			helpDelete(node, "help concurrent")
			continue
		}

		// 1 pop step
		if node.next.compareAndSet(nodeNext, nodeNext, false, true) {
			// 2, 3 step
			helpDelete(node, "1st step")
			next := node.next.getReference()
			// 4 step
			helpInsert(prev, next, "popLeft New")

			return node.Value
		}
	}
}

// PopRight returns the last element of list l or nil if the list is empty.
func (l *List) PopRight() interface{} {
	next := l.tail
	node := next.prev.getReference()

	for {
		if !node.next.compareAndSet(next, next, false, false) {
			node = helpInsert(node, next, "popRight")
			continue
		}

		if node == l.head {
			return nil
		}

		if node.next.compareAndSet(next, next, false, true) {
			helpDelete(node, "")
			prev := node.prev.getReference()
			helpInsert(prev, next, "popRight")

			return node.Value
		}
	}
}

// PushLeft inserts a new element e with value v at the front of list l and returns e.
func (l *List) PushLeft(v interface{}) *Element {
	node := &Element{Value: v}
	prev := l.head
	next := prev.next.getReference()

	for {
		if !prev.next.compareAndSet(next, next, false, false) {
			next = prev.next.getReference()
			continue
		}

		node.prev = newAtomicMarkableReference(prev, false)
		node.next = newAtomicMarkableReference(next, false)

		if prev.next.compareAndSet(next, node, false, false) {
			break
		}
	}

	pushCommon(node, next)

	return node
}

// PushRight inserts a new element e with value v at the back of list l and returns e.
func (l *List) PushRight(v interface{}) *Element {
	node := &Element{Value: v}
	next := l.tail
	prev := next.prev.getReference()

	for {
		if !prev.next.compareAndSet(next, next, false, false) {
			// concurrent push inserted -> get new prev
			prev = helpInsert(prev, next, "concurrentPushRight")
			continue
		}

		// 0 push step
		node.prev = newAtomicMarkableReference(prev, false)
		node.next = newAtomicMarkableReference(next, false)
		// 1 push step
		if prev.next.compareAndSet(next, node, false, false) {
			break
		}
	}

	// 2 push step
	pushCommon(node, next)

	return node
}

func pushCommon(node, next *Element) {
	for {
		link1 := next.prev
		if link1.isMarked() || !node.next.compareAndSet(next, next, false, false) {
			break
		}

		if next.prev.compareAndSet(link1.getReference(), node, false, false) {
			if node.prev.isMarked() {
				helpInsert(node, next, "pushCommon")
			}
			break
		}
	}
}

// Next returns the next list element or nil.
func (l *List) Next(node *Element) *Element {
	for node != l.tail {
		if node == nil {
			break
		}

		next := node.next.getReference()
		if next == nil {
			break
		}

		removed, nextNext := next.next.get()
		if removed {
			// The next pointer of the node behind me has the deleted mark set
			removed2, nodeNext := node.next.get()
			if !removed2 || nodeNext != next {
				setMark(next.prev)
				node.next.compareAndSet(next, nextNext, false, false) // next removed == false?
				continue
			}
		}

		node = next

		if !removed {
			return next
		}
	}

	return nil
}

// Prev returns the previous list element or nil.
func (l *List) Prev(node *Element) *Element {
	for node != l.head {
		if node == nil {
			break
		}

		prev := node.prev.getReference()
		if prev == nil {
			break
		}

		prevNext := prev.next.getReference()
		removed := node.next.isMarked()

		if prevNext == node && !removed {
			return prev
		} else if removed {
			node = l.Next(node)
		} else {
			prev = l.correctPrev(prev, node)
		}
	}

	return nil
}

/**
 * Correct node.prev to the closest previous node
 * helpInsert is very weak - does not reset node.prev to the actual prev.next
 * but just tries to set node.prev to the given suggestion of a prev node
 * (for 2 push step, 4 pop step)
 */
func helpInsert(prev *Element, node *Element, method string) *Element {
	// last = is the last node : last.next == prev and it is not marked as removed
	var last, nodePrev *Element

	for {
		removed, prevNext := prev.next.get()

		if removed {
			if last != nil {
				markPrev(prev)
				next2 := prev.next.getReference()
				last.next.compareAndSet(prev, next2, false, false)
				prev = last
				last = nil
			} else {
				prevNext = prev.prev.getReference()
				prev = prevNext
			}
			continue
		}

		removed, nodePrev = node.prev.get()
		if removed {
			break
		}

		// prev is not the previous node of node
		if prevNext != node {
			last = prev
			prev = prevNext
			continue
		}

		if nodePrev == prev {
			break
		}

		if prev.next.getReference() == node && node.prev.compareAndSet(nodePrev, prev, false, false) {
			if prev.prev.isMarked() {
				continue
			}
			break
		}
	}

	return prev
}

// 2 and 3 pop steps
func helpDelete(node *Element, place string) {
	markPrev(node)

	prev := node.prev.getReference()
	next := node.next.getReference()
	var last *Element

	for {
		if prev == next {
			break
		}

		if next.next.isMarked() {
			markPrev(next)
			next = next.next.getReference()
			continue
		}

		removed, prevNext := prev.next.get()
		if removed {
			if last != nil {
				markPrev(prev)
				next2 := prev.next.getReference()
				last.next.compareAndSet(prev, next2, false, false)
				prev = last
				last = nil
			} else {
				prevNext = prev.prev.getReference()
				prev = prevNext
				// assert(prev != nil)
			}
			continue
		}

		if prevNext != node {
			last = prev
			prev = prevNext
			continue
		}

		if prev.next.compareAndSet(node, next, false, false) {
			break
		}
	}
}

func markPrev(node *Element) {
	for {
		link1 := node.prev
		if link1.isMarked() ||
			node.prev.compareAndSet(link1.getReference(), link1.getReference(), false, true) {
			break
		}
	}
}
