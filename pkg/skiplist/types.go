package skiplist

import (
	"math/rand"
	"sync"
)

type (
	elementNode struct {
		forward []*Element
	}

	// Element represents an element in the list.
	Element struct {
		elementNode
		key   string
		value int64
	}

	// SkipList represents the list.
	SkipList struct {
		elementNode
		randSource       rand.Source
		searchNodesCache []*elementNode
		probability      float64
		probTable        []float64
		maxLevel         int
		length           int
		mutex            sync.RWMutex
	}
)

// Key returns the given Element key.
func (e *Element) Key() string {
	return e.key
}

// Value returns the given Element value.
func (e *Element) Value() int64 {
	return e.value
}

// Next returns the adjacent next Element if existed or nil otherwise.
func (e *Element) Next() *Element {
	return e.forward[0]
}

// NextLevel returns the adjacent next Element at provided level if existed or nil otherwise.
func (e *Element) NextLevel(level int) *Element {
	if level >= len(e.forward) || level < 0 {
		return nil
	}

	return e.forward[level]
}
