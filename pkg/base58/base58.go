// Package base58 implements Base58 Encoder interface.
package base58

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	source = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	ln     = uint64(58)
)

var (
	encodeMap [58]byte
	decodeMap [256]int
	// ErrEmptyString for empty Base58 encoded string.
	ErrEmptyString = errors.New("base58: empty Base58 encoded string")
)

func init() {
	for i := range decodeMap {
		decodeMap[i] = -1
	}
	for i := range source {
		encodeMap[i] = source[i]
		decodeMap[encodeMap[i]] = i
	}
}

// Encode returns Base58 encoded string.
func Encode(n uint64) string {
	if n == 0 {
		return string(encodeMap[:1])
	}

	b, div := make([]byte, 0, binary.MaxVarintLen64), uint64(0)
	for n > 0 {
		div = n / ln
		b = append(b, encodeMap[n-div*ln])
		n = div
	}
	Reverse(b)
	return string(b)
}

// Decode returns Base58 decoded unsigned int64.
func Decode(s string) (uint64, error) {
	if s == "" {
		return 0, ErrEmptyString
	}

	var (
		n uint64
		c int
	)
	for i := range s {
		c = decodeMap[s[i]]
		if c < 0 {
			return 0, fmt.Errorf("base58: invalid character <%s>", string(s[i]))
		}
		n = n*58 + uint64(c)
	}
	return n, nil
}

// Reverse your byte slice.
func Reverse(b []byte) {
	for l, r := 0, len(b)-1; l < r; l, r = l+1, r-1 {
		b[l], b[r] = b[r], b[l]
	}
}
