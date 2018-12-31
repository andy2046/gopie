// Package drf implements Dominant Resource Fairness.
package drf

import (
	"container/heap"
	"sync"
)

// New create a DRF Cluster.
func New(clusterResource map[Typ]float64, clusterNodes ...*Node) (DRF, error) {
	if len(clusterResource) == 0 {
		return EmptyDRF, ErrEmptyResource
	}
	if len(clusterNodes) == 0 {
		return EmptyDRF, ErrEmptyNodes
	}

	drf := DRF{
		clusterResource:  clusterResource,
		consumedResource: make(map[Typ]float64),
		nodes:            make(nodeQueue, len(clusterNodes)),
		mu:               &sync.RWMutex{},
	}
	i := 0
	for _, n := range clusterNodes {
		n.mu.Lock()
		n.index = i
		if n.allocated == nil {
			n.allocated = make(map[Typ]float64)
		}
		if n.demand == nil {
			n.demand = make(map[Typ]float64)
		}
		drf.nodes[i] = n
		n.mu.Unlock()
		i++
	}
	if drf.clusterResource == nil {
		drf.clusterResource = make(map[Typ]float64)
	}
	heap.Init(&drf.nodes)
	return drf, nil
}

// NewNode create a Cluster Node.
func NewNode(demand ...map[Typ]float64) *Node {
	n := Node{
		allocated: make(map[Typ]float64),
	}
	if len(demand) > 0 {
		n.demand = demand[0]
	}
	if n.demand == nil {
		n.demand = make(map[Typ]float64)
	}
	return &n
}

// UpdateDemand add delta to existing demand.
func (n *Node) UpdateDemand(delta map[Typ]float64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for k, v := range delta {
		n.demand[k] = n.demand[k] + v
	}
}

// NextTask run next task with lowest dominant share.
func (drf DRF) NextTask() error {
	drf.mu.Lock()
	defer drf.mu.Unlock()
	if len(drf.nodes) == 0 {
		return ErrEmptyNodes
	}
	n := drf.nodes[0]
	n.mu.Lock()
	defer n.mu.Unlock()
	if drf.checkIfResourceUsageOverLimit(n) {
		return ErrResourceSaturated
	}
	n.updateAllocated(n.demand)
	drf.updateConsumed(n.demand)
	drf.computeDShare(n)
	heap.Fix(&drf.nodes, n.index)
	return nil
}

// UpdateResource add delta to Cluster Resource.
func (drf DRF) UpdateResource(delta map[Typ]float64) {
	drf.mu.Lock()
	defer drf.mu.Unlock()
	for k, v := range delta {
		drf.clusterResource[k] = drf.clusterResource[k] + v
	}
}

// AddNode add new Node to DRF Cluster.
func (drf DRF) AddNode(n *Node) {
	drf.mu.Lock()
	defer drf.mu.Unlock()
	n.mu.Lock()
	defer n.mu.Unlock()
	n.allocated = make(map[Typ]float64)
	n.dShare = 0
	if n.demand == nil {
		// TODO: error out if demand is empty
		n.demand = make(map[Typ]float64)
	}
	drf.nodes.Push(n)
}

// RemoveNode remove Node from DRF Cluster.
func (drf DRF) RemoveNode(n *Node) {
	drf.mu.Lock()
	defer drf.mu.Unlock()
	n.mu.Lock()
	defer n.mu.Unlock()
	heap.Remove(&drf.nodes, n.index)
	n.index = -1
	for k, v := range n.allocated {
		drf.consumedResource[k] = drf.consumedResource[k] - v
	}
	n.allocated = nil
	n.dShare = 0
}

// Allocated return all the allocated resource for node.
func (n *Node) Allocated() map[Typ]float64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	allc := make(map[Typ]float64, len(n.allocated))
	for k, v := range n.allocated {
		allc[k] = allc[k] + v
	}
	return allc
}

// Consumed return all the consumed resource by cluster.
func (drf DRF) Consumed() map[Typ]float64 {
	drf.mu.RLock()
	defer drf.mu.RUnlock()
	consm := make(map[Typ]float64, len(drf.consumedResource))
	for k, v := range drf.consumedResource {
		consm[k] = consm[k] + v
	}
	return consm
}

// Resource return all the cluster resource.
func (drf DRF) Resource() map[Typ]float64 {
	drf.mu.RLock()
	defer drf.mu.RUnlock()
	cr := make(map[Typ]float64, len(drf.clusterResource))
	for k, v := range drf.clusterResource {
		cr[k] = cr[k] + v
	}
	return cr
}

func (drf DRF) checkIfResourceUsageOverLimit(n *Node) bool {
	for k := range n.demand {
		if drf.consumedResource[k]+n.demand[k] > drf.clusterResource[k] {
			return true
		}
	}
	return false
}

func (drf DRF) computeDShare(n *Node) {
	temp := n.dShare
	for k := range n.allocated {
		if r, ok := drf.clusterResource[k]; ok && r > 0 {
			if n.allocated[k]/r > temp {
				temp = n.allocated[k] / r
			}
		}
	}
	n.dShare = temp
}

func (n *Node) updateAllocated(demand map[Typ]float64) {
	for k := range demand {
		n.allocated[k] = n.allocated[k] + demand[k]
	}
}

func (drf DRF) updateConsumed(demand map[Typ]float64) {
	for k := range demand {
		drf.consumedResource[k] = drf.consumedResource[k] + demand[k]
	}
}
