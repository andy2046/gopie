// Package singleton provides a singleton implementation.
package singleton

import "sync"

// Instance is the singleton instance.
type Instance struct {
	Values map[interface{}]interface{}
}

var (
	once     sync.Once
	instance *Instance
)

// New returns the singleton instance.
func New() *Instance {
	once.Do(func() {
		instance = &Instance{
			Values: make(map[interface{}]interface{}),
		}
	})
	return instance
}
