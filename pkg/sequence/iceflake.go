package sequence

import (
	"time"
)

// Iceflake is the interface for snowflake similar sequence generator.
type Iceflake interface {
	Sequencer
	// StartTime defines the time since which
	// the Iceflake time is defined as the elapsed time.
	StartTime() time.Time
	// BitLenSequence defines the bit length of sequence number,
	// and the bit length of time is 63 - BitLenSequence().
	BitLenSequence() uint8
}
