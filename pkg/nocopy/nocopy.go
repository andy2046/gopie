// Package nocopy implements the interface for -copylocks checker from `go vet`.
package nocopy

// NoCopy can be embedded into structs which must not be copied
// after the first use.
type NoCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*NoCopy) Lock() {}

// Unlock is not required by -copylocks checker from `go vet`.
func (*NoCopy) Unlock() {}
