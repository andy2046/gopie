package drf

import (
	"container/heap"
	"fmt"
	"sync"
)

type (
	// Typ represents resource type.
	Typ string

	// Node represents a Cluster Node.
	Node struct {
		index     int
		allocated map[Typ]float64
		demand    map[Typ]float64
		dShare    float64 // Dominant Share
		mu        sync.RWMutex
	}

	// A nodeQueue implements heap.Interface and holds Nodes.
	nodeQueue []*Node

	// DRF represents a DRF Cluster.
	DRF struct {
		clusterResource  map[Typ]float64
		consumedResource map[Typ]float64
		nodes            nodeQueue
		mu               *sync.RWMutex
	}
)

const (
	// CPU resource.
	CPU = Typ("CPU")

	// MEMORY resource.
	MEMORY = Typ("MEMORY")
)

var (
	// ErrResourceSaturated when there is not enough resource to run next task.
	ErrResourceSaturated = fmt.Errorf("Fatal: resource has been saturated")

	// ErrEmptyResource when cluster resource is empty.
	ErrEmptyResource = fmt.Errorf("Fatal: empty cluster resource")

	// ErrEmptyNodes when cluster nodes is empty.
	ErrEmptyNodes = fmt.Errorf("Fatal: empty cluster nodes")

	// EmptyDRF is empty DRF.
	EmptyDRF = DRF{}
)

func (nq nodeQueue) Len() int { return len(nq) }

func (nq nodeQueue) Less(i, j int) bool {
	// Pop the lowest dShare
	return nq[i].dShare < nq[j].dShare
}

func (nq nodeQueue) Swap(i, j int) {
	nq[i], nq[j] = nq[j], nq[i]
	nq[i].index = i
	nq[j].index = j
}

func (nq *nodeQueue) Push(x interface{}) {
	n := len(*nq)
	item := x.(*Node)
	item.index = n
	*nq = append(*nq, item)
}

func (nq *nodeQueue) Pop() interface{} {
	old := *nq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*nq = old[0 : n-1]
	return item
}

// update modifies an Node in the queue.
func (nq *nodeQueue) update(item *Node, dShare float64) {
	item.dShare = dShare
	heap.Fix(nq, item.index)
}
