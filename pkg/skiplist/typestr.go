// +build ignore

package skiplist

import (
	"math/rand"
	"sync"
)

type (
	elementNodeS struct {
		forward []*ElementS
	}

	// ElementS represents an element in the list.
	ElementS struct {
		elementNodeS
		key   string
		value string
	}

	// S represents the list.
	S struct {
		elementNodeS
		randSource       rand.Source
		searchNodesCache []*elementNodeS
		probability      float64
		probTable        []float64
		maxLevel         int
		length           int
		mutex            sync.RWMutex
	}
)

// Key returns the given Element key.
func (e *ElementS) Key() string {
	return e.key
}

// Value returns the given Element value.
func (e *ElementS) Value() string {
	return e.value
}

// Next returns the adjacent next Element if existed or nil otherwise.
func (e *ElementS) Next() *ElementS {
	return e.forward[0]
}

// NextLevel returns the adjacent next Element at provided level if existed or nil otherwise.
func (e *ElementS) NextLevel(level int) *ElementS {
	if level >= len(e.forward) || level < 0 {
		return nil
	}

	return e.forward[level]
}
