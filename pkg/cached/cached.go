// Package cached implements Caching Devil pattern.
package cached

import "time"

type (
	// Lease is a per-cache-key lock preventing thundering herds and stale sets.
	Lease interface {
		// Nonce is unique for each Lease.
		Nonce() string
	}

	// Lessor is the Lease helper.
	Lessor interface {
		// NewLease create a new Lease.
		NewLease() Lease

		// FromValue construct a Lease from nonce value.
		FromValue(nonce []byte) (Lease, error)

		// MustFromValue is like FromValue but panics if nonce value is not a Lease.
		MustFromValue(nonce []byte) Lease

		// IsLease returns true if value is a Lease, false otherwise.
		IsLease(value []byte) bool
	}

	// Cache represents the Caching layer.
	Cache interface {
		// Get returns a list of []byte representing the values associated with the provided keys.
		// The value associated with a certain key will be in the same position in the returned list
		// as the key is in the keys list.
		Get(keys ...string) [][]byte

		// Set associates the provided value with the provided key in the cache layer.
		Set(key string, value []byte) error

		// AtomicAdd set the provided value for key if and only if the key has not already been set.
		// returns true if it succeeds, false otherwise.
		AtomicAdd(key string, value []byte) bool

		// AtomicCheckAndSet set the valueToSet for the provided key if and only if the key is currently associated with expectedValue.
		// returns true if it succeeds, false otherwise.
		AtomicCheckAndSet(key string, expectedValue, valueToSet []byte) bool
	}

	// TruthTeller is the function to fetch the value associated with the looked up key
	// from the source of truth data store.
	TruthTeller func(key string) []byte

	// Cacher manages Caching Devil.
	Cacher struct {
		interval time.Duration
		cache    Cache
		lessor   Lessor
		teller   TruthTeller
	}
)

// New returns a new Cacher.
// interval is the retry waiting time if another request is holding the lease.
func New(interval time.Duration, cache Cache, lessor Lessor, teller TruthTeller) Cacher {
	return Cacher{
		interval: interval,
		cache:    cache,
		lessor:   lessor,
		teller:   teller,
	}
}

// Read try to retrieve the value associated with the provided key.
func (c Cacher) Read(key string) []byte {
	return c.read(key, []Lease{})
}

func (c Cacher) read(key string, previouslySeenLeases []Lease) []byte {
	cacheKeysToLookUp := []string{key}
	for _, lease := range previouslySeenLeases {
		cacheKeysToLookUp = append(cacheKeysToLookUp, lease.Nonce())
	}

	// looking up the provided key as well as all previously seen leases
	var (
		valueForKey     []byte
		valuesFromCache = c.cache.Get(cacheKeysToLookUp...)
	)

	if len(valuesFromCache) == 0 {
		return nil
	}
	valueForKey, valuesFromCache = valuesFromCache[0], valuesFromCache[1:]

	// check if the value is stashed behind one of the leases we've previously seen
	// avoid reader to experience a ridiculous amount of latency as it waits for other readers to populate the cache key
	for _, valueForPreviouslySeenLease := range valuesFromCache {
		if valueForPreviouslySeenLease != nil {
			return valueForPreviouslySeenLease
		}
	}

	if valueForKey == nil {
		// value is not in the cache
		newLease := c.lessor.NewLease()
		nonceBytes := []byte(newLease.Nonce())
		leaseAdded := c.cache.AtomicAdd(key, nonceBytes)

		if leaseAdded {
			// managed to acquire a lease on this key,
			// now populate the cache with the value from data store
			valueFromTruth := c.teller(key)

			// avoid cache poisoning with stale value
			// if CAS returns false, the value is invalidated by write
			_ = c.cache.AtomicCheckAndSet(key, nonceBytes, valueFromTruth)

			// avoid reader to experience a ridiculous amount of latency as it waits for other readers to populate the cache key
			c.cache.Set(newLease.Nonce(), valueFromTruth)
			return valueFromTruth
		}

		// another request managed to acquire the lease before me, retry
		return c.read(key, previouslySeenLeases)
	} else if c.lessor.IsLease(valueForKey) {
		// another request is holding the lease on this key, try again later
		time.Sleep(c.interval)

		previouslySeenLeases = append(previouslySeenLeases, c.lessor.MustFromValue(valueForKey))
		return c.read(key, previouslySeenLeases)
	} else {
		// got the value from cache, return it
		return valueForKey
	}
}
