// Package bitflag implements bit flag.
package bitflag

// Flag is an 8 bits flag.
type Flag byte

const siz uint8 = 7

// SetAll set all the flags.
func (f *Flag) SetAll(opts ...Flag) {
	for _, o := range opts {
		*f |= o
	}
}

// ToggleAll toggle (XOR) all the flags.
func (f *Flag) ToggleAll(opts ...Flag) {
	for _, o := range opts {
		*f ^= o
	}
}

// ClearAll clear all the flags.
func (f *Flag) ClearAll(opts ...Flag) {
	for _, o := range opts {
		*f &^= o
	}
}

// AreAllSet check if all the flags are set.
func (f Flag) AreAllSet(opts ...Flag) bool {
	for _, o := range opts {
		if f&o == 0 {
			return false
		}
	}
	return true
}

// IsAnySet check if any one flag is set.
func (f Flag) IsAnySet(opts ...Flag) bool {
	for _, o := range opts {
		if f&o > 0 {
			return true
		}
	}
	return false
}

// IsSet check if the bit at `n` is set.
// `n` should be less than `8`.
func (f Flag) IsSet(n uint8) bool {
	if n > siz || f&(1<<n) == 0 {
		return false
	}
	return true
}

// Set set a single bit at `n`.
// `n` should be less than `8`.
func (f *Flag) Set(n uint8) {
	if n > siz {
		return
	}
	*f |= (1 << n)
}

// Toggle toggle (XOR) a single bit at `n`.
// `n` should be less than `8`.
func (f *Flag) Toggle(n uint8) {
	if n > siz {
		return
	}
	*f ^= (1 << n)
}

// Clear clear a single bit at `n`.
// `n` should be less than `8`.
func (f *Flag) Clear(n uint8) {
	if n > siz {
		return
	}
	*f &^= (1 << n)
}

// Reset reset the flag.
func (f *Flag) Reset() {
	*f = 0
}
