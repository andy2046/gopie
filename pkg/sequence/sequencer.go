// Package sequence implements Iceflake sequence generator interface.
// ```bash
// Iceflake is the interface for snowflake similar sequence generator.
//
// Iceflake algorithm:
//
// +-------+--------------------+----------+
// | sign  | delta milliseconds | sequence |
// +-------+--------------------+----------+
// | 1 bit | 63-n bits          | n bits   |
//
// sequence (n bits)
// The last custom n bits, represents sequence within the one millisecond.
//
// delta milliseconds (63-n bits)
// The next 63-n bits, represents delta milliseconds since a custom epoch.
// ```
package sequence

// Sequencer is the interface for sequence generator.
type Sequencer interface {
	// Next returns the next sequence.
	Next() (uint64, error)
	// NextN reserves the next `n` sequences and returns the first one,
	// `n` should not be less than 1.
	NextN(n int) (uint64, error)
	// MachineID returns the unique ID of the instance.
	MachineID() uint64
}
