// +build ignore

package skiplist

import (
	"math/rand"
	"time"
)

// NewS creates a new skip list with string type value and provided maxLevel or defaultMaxLevel.
func NewS(maxLevel ...int) *S {
	level := defaultMaxLevel
	if len(maxLevel) > 0 && maxLevel[0] >= 1 && maxLevel[0] <= 64 {
		level = maxLevel[0]
	}

	return &S{
		elementNodeS:     elementNodeS{forward: make([]*ElementS, level)},
		searchNodesCache: make([]*elementNodeS, level),
		maxLevel:         level,
		randSource:       rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:      defaultProbability,
		probTable:        probabilityTable(defaultProbability, level),
	}
}

// Front returns the first element in the list.
func (list *S) Front() *ElementS {
	return list.forward[0]
}

// Len returns the list length.
func (list *S) Len() int {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	return list.length
}

// MaxLevel sets current max level of skip list and returns the previous max level.
// If `level` < 1, it does not set current max level.
func (list *S) MaxLevel(level int) int {
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
		f := make([]*ElementS, level)
		copy(f, list.forward)
		for i := range list.forward {
			list.forward[i] = nil // avoid mem leak
		}
		list.forward = f
		for i := range list.searchNodesCache {
			list.searchNodesCache[i] = nil // avoid mem leak
		}
		list.searchNodesCache = make([]*elementNodeS, level)
		list.probTable = probabilityTable(list.probability, level)
	}

	return prev
}

// Set upserts the value with provided key into the list and returns the element.
func (list *S) Set(key string, value string) *ElementS {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var element *ElementS
	fingerSearches := list.fingerSearchNodes(key)

	if element = fingerSearches[0].forward[0]; element != nil && element.key <= key {
		element.value = value
		return element
	}

	element = &ElementS{
		elementNodeS: elementNodeS{
			forward: make([]*ElementS, list.randLevel()),
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
func (list *S) Get(key string) *ElementS {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var next *ElementS
	prev := &list.elementNodeS

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.forward[i]

		for next != nil && key > next.key {
			prev = &next.elementNodeS
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
func (list *S) Remove(key string) *ElementS {
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
func (list *S) fingerSearchNodes(key string) []*elementNodeS {
	var next *ElementS
	prev := &list.elementNodeS
	fingerSearches := list.searchNodesCache

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.forward[i]

		for next != nil && key > next.key {
			prev = &next.elementNodeS
			next = next.forward[i]
		}

		fingerSearches[i] = prev
	}

	return fingerSearches
}

func (list *S) randLevel() int {
	r := float64(list.randSource.Int63()) / (1 << 63)
	level := 1
	for level < list.maxLevel && r < list.probTable[level] {
		level++
	}
	return level
}

// Probability sets the P value of skip list and returns the previous P.
// If `newProbability` < 0, it does not set current P.
func (list *S) Probability(newProbability float64) float64 {
	p := list.probability
	if newProbability >= 0 {
		list.mutex.Lock()
		defer list.mutex.Unlock()
		list.probability = newProbability
		list.probTable = probabilityTable(newProbability, list.maxLevel)
	}
	return p
}
