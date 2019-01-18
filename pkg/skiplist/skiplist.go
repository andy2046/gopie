// Package skiplist provides a Skip List implementation.
package skiplist

import (
	"math"
	"math/rand"
	"time"
)

const (
	defaultMaxLevel int = 16

	defaultProbability float64 = 1 / math.E
)

// New creates a new skip list with provided maxLevel or defaultMaxLevel.
func New(maxLevel ...int) *SkipList {
	level := defaultMaxLevel
	if len(maxLevel) > 0 && maxLevel[0] >= 1 && maxLevel[0] <= 64 {
		level = maxLevel[0]
	}

	return &SkipList{
		elementNode:      elementNode{forward: make([]*Element, level)},
		searchNodesCache: make([]*elementNode, level),
		maxLevel:         level,
		randSource:       rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:      defaultProbability,
		probTable:        probabilityTable(defaultProbability, level),
	}
}

// Front returns the first element in the list.
func (list *SkipList) Front() *Element {
	return list.forward[0]
}

// Len returns the list length.
func (list *SkipList) Len() int {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	return list.length
}

// MaxLevel sets current max level of skip list and returns the previous max level.
// If `level` < 1, it does not set current max level.
func (list *SkipList) MaxLevel(level int) int {
	if level < 1 || list.maxLevel == level {
		return list.maxLevel
	}

	list.mutex.Lock()
	defer list.mutex.Unlock()
	prev := list.maxLevel
	list.maxLevel = level

	switch {
	case prev > level:
		for k, n := level, len(list.forward); k < n; k++ {
			list.forward[k] = nil // avoid mem leak
		}
		list.forward = list.forward[:level]
		for k, n := level, len(list.searchNodesCache); k < n; k++ {
			list.searchNodesCache[k] = nil // avoid mem leak
		}
		list.searchNodesCache = list.searchNodesCache[:level]
		list.probTable = list.probTable[:level]
	case prev < level:
		f := make([]*Element, level)
		copy(f, list.forward)
		for i := range list.forward {
			list.forward[i] = nil // avoid mem leak
		}
		list.forward = f
		for i := range list.searchNodesCache {
			list.searchNodesCache[i] = nil // avoid mem leak
		}
		list.searchNodesCache = make([]*elementNode, level)
		list.probTable = probabilityTable(list.probability, level)
	}

	return prev
}

// Set upserts the value with provided key into the list and returns the element.
func (list *SkipList) Set(key string, value int64) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var element *Element
	fingerSearches := list.fingerSearchNodes(key)

	if element = fingerSearches[0].forward[0]; element != nil && element.key <= key {
		element.value = value
		return element
	}

	element = &Element{
		elementNode: elementNode{
			forward: make([]*Element, list.randLevel()),
		},
		key:   key,
		value: value,
	}

	for i := range element.forward {
		element.forward[i] = fingerSearches[i].forward[i]
		fingerSearches[i].forward[i] = element
	}

	list.length++
	return element
}

// Get searches element by provided key and returns the element if found or nil otherwise.
func (list *SkipList) Get(key string) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var next *Element
	prev := &list.elementNode

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.forward[i]

		for next != nil && key > next.key {
			prev = &next.elementNode
			next = next.forward[i]
		}
	}

	if next != nil && next.key <= key {
		return next
	}

	return nil
}

// Remove deletes the element with provided key from the list,
// and returns the removed element if found or nil otherwise.
func (list *SkipList) Remove(key string) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	fingerSearches := list.fingerSearchNodes(key)

	if element := fingerSearches[0].forward[0]; element != nil && element.key <= key {
		for k, v := range element.forward {
			fingerSearches[k].forward[k] = v
		}

		list.length--
		return element
	}

	return nil
}

// fingerSearchNodes returns a list of nodes, where nodes[i] contains a pointer to the rightmost node
// of level i or higher that is to the left of the location of the `key`.
func (list *SkipList) fingerSearchNodes(key string) []*elementNode {
	var next *Element
	prev := &list.elementNode
	fingerSearches := list.searchNodesCache

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.forward[i]

		for next != nil && key > next.key {
			prev = &next.elementNode
			next = next.forward[i]
		}

		fingerSearches[i] = prev
	}

	return fingerSearches
}

func (list *SkipList) randLevel() int {
	r := float64(list.randSource.Int63()) / (1 << 63)
	level := 1
	for level < list.maxLevel && r < list.probTable[level] {
		level++
	}
	return level
}

// Probability sets the P value of skip list and returns the previous P.
// If `newProbability` < 0, it does not set current P.
func (list *SkipList) Probability(newProbability float64) float64 {
	p := list.probability
	if newProbability >= 0 {
		list.mutex.Lock()
		defer list.mutex.Unlock()
		list.probability = newProbability
		list.probTable = probabilityTable(newProbability, list.maxLevel)
	}
	return p
}

// probabilityTable stores the probability for a new node appearing in a given level,
// probability is in [0, 1], maxLevel is in (0, 64].
func probabilityTable(probability float64, maxLevel int) []float64 {
	t := make([]float64, 0, maxLevel)
	for i := 1; i <= maxLevel; i++ {
		t = append(t, math.Pow(probability, float64(i-1)))
	}
	return t
}
