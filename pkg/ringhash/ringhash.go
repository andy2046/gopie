// Package ringhash provides a ring hash implementation.
package ringhash

import (
	"errors"
	"hash/crc64"
	"io"
	"math"
	"sort"
	"strconv"
	"sync"
)

type (
	// Hash is the hash function.
	Hash func(key string) uint64

	// Node is the node in the ring.
	Node struct {
		Name string
		Load int64
	}

	// Ring is the data store for keys hash map.
	Ring struct {
		hashFn          Hash
		replicas        int
		balancingFactor float64
		hashKeyMap      map[uint64]string
		hashes          []uint64
		keyLoadMap      map[string]*Node
		totalLoad       int64

		mu sync.RWMutex
	}

	// Config is the config for hash ring.
	Config struct {
		HashFn          Hash
		Replicas        int
		BalancingFactor float64
	}

	// Option applies config to Config.
	Option = func(*Config) error
)

var (
	// hashCRC64 uses the 64-bit Cyclic Redundancy Check (CRC-64) with the ECMA polynomial.
	hashCRC64 = crc64.New(crc64.MakeTable(crc64.ECMA))
	// ErrNoNode when there is no node added into the hash ring.
	ErrNoNode = errors.New("no node added")
	// ErrNodeNotFound when no node found in LoadMap.
	ErrNodeNotFound = errors.New("node not found in LoadMap")
	// DefaultConfig is the default config for hash ring.
	DefaultConfig = Config{
		HashFn:          hash,
		Replicas:        10,
		BalancingFactor: 1.25,
	}
)

func hash(key string) uint64 {
	hashCRC64.Reset()
	_, err := io.WriteString(hashCRC64, key)
	if err != nil {
		panic(err)
	}
	return hashCRC64.Sum64()
}

// New returns a new Ring.
func New(options ...Option) *Ring {
	c := DefaultConfig
	setOption(&c, options...)
	r := &Ring{
		replicas:        c.Replicas,
		balancingFactor: c.BalancingFactor,
		hashFn:          c.HashFn,
		hashKeyMap:      make(map[uint64]string),
		hashes:          []uint64{},
		keyLoadMap:      map[string]*Node{},
	}
	return r
}

// IsEmpty returns true if there is no node in the ring.
func (r *Ring) IsEmpty() bool {
	return len(r.hashKeyMap) == 0
}

// AddNode adds Node with key as name to the hash ring.
func (r *Ring) AddNode(keys ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, key := range keys {
		if _, ok := r.keyLoadMap[key]; ok {
			continue
		}

		r.keyLoadMap[key] = &Node{Name: key, Load: 0}

		for i := 0; i < r.replicas; i++ {
			h := r.hashFn(key + strconv.Itoa(i))
			r.hashes = append(r.hashes, h)
			r.hashKeyMap[h] = key
		}
	}

	// sort hashes ascendingly
	sort.Slice(r.hashes, func(i, j int) bool {
		if r.hashes[i] < r.hashes[j] {
			return true
		}
		return false
	})
}

// GetNode returns the closest node in the hash ring to the provided key.
func (r *Ring) GetNode(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.IsEmpty() {
		return "", ErrNoNode
	}

	h := r.hashFn(key)
	idx := r.search(h)
	return r.hashKeyMap[r.hashes[idx]], nil
}

// GetLeastNode uses consistent hashing with bounded loads to get the least loaded node.
func (r *Ring) GetLeastNode(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.IsEmpty() {
		return "", ErrNoNode
	}

	h := r.hashFn(key)
	idx := r.search(h)

	i := idx
	var count int
	for {
		count++
		if count >= len(r.hashes) {
			panic("not enough space to distribute load")
		}
		node := r.hashKeyMap[r.hashes[i]]
		if r.loadOK(node) {
			return node, nil
		}
		i++
		if i >= len(r.hashKeyMap) {
			i = 0
		}
	}
}

// UpdateLoad sets load of the given node to the given load.
func (r *Ring) UpdateLoad(node string, load int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.keyLoadMap[node]; !ok {
		return
	}
	r.totalLoad -= r.keyLoadMap[node].Load
	r.keyLoadMap[node].Load = load
	r.totalLoad += load
}

// Add increases load of the given node by 1,
// should only be used with GetLeast.
func (r *Ring) Add(node string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.keyLoadMap[node]; !ok {
		return false
	}
	r.keyLoadMap[node].Load++
	r.totalLoad++
	return true
}

// Done decreases load of the given node by 1,
// should only be used with GetLeast.
func (r *Ring) Done(node string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.keyLoadMap[node]; !ok {
		return false
	}
	r.keyLoadMap[node].Load--
	r.totalLoad--
	return true
}

// RemoveNode deletes node from the hash ring.
func (r *Ring) RemoveNode(node string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.keyLoadMap[node]; !ok {
		return false
	}

	for i := 0; i < r.replicas; i++ {
		h := r.hashFn(node + strconv.Itoa(i))
		delete(r.hashKeyMap, h)
		r.removeFromHashes(h)
	}
	delete(r.keyLoadMap, node)
	return true
}

// Nodes returns the list of nodes in the hash ring.
func (r *Ring) Nodes() (nodes []string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for k := range r.keyLoadMap {
		nodes = append(nodes, k)
	}
	return
}

// Loads returns the loads of all the nodes in the hash ring.
func (r *Ring) Loads() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loads := map[string]int64{}
	for k, v := range r.keyLoadMap {
		loads[k] = v.Load
	}
	return loads
}

// MaxLoad returns the maximum load for a single node in the hash ring,
// which is (totalLoad/numberOfNodes)*balancingFactor.
func (r *Ring) MaxLoad() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.totalLoad == 0 {
		r.totalLoad = 1
	}
	var avgLoadPerNode float64
	avgLoadPerNode = float64(r.totalLoad / int64(len(r.keyLoadMap)))
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	avgLoadPerNode = math.Ceil(avgLoadPerNode * r.balancingFactor)
	return int64(avgLoadPerNode)
}

func (r *Ring) removeFromHashes(h uint64) {
	for i := 0; i < len(r.hashes); i++ {
		if r.hashes[i] == h {
			r.hashes = append(r.hashes[:i], r.hashes[i+1:]...)
		}
	}
}

func setOption(c *Config, options ...func(*Config) error) error {
	for _, opt := range options {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}

func (r *Ring) search(key uint64) int {
	idx := sort.Search(len(r.hashes), func(i int) bool {
		return r.hashes[i] >= key
	})

	if idx >= len(r.hashes) {
		idx = 0
	}
	return idx
}

func (r *Ring) loadOK(key string) bool {
	if r.totalLoad < 0 {
		r.totalLoad = 0
	}

	var avgLoadPerNode float64
	avgLoadPerNode = float64((r.totalLoad + 1) / int64(len(r.keyLoadMap)))
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	avgLoadPerNode = math.Ceil(avgLoadPerNode * r.balancingFactor)

	node, ok := r.keyLoadMap[key]
	if !ok {
		panic(ErrNodeNotFound)
	}

	if float64(node.Load)+1 <= avgLoadPerNode {
		return true
	}

	return false
}
