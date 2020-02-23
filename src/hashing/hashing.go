package hashing

import (
	"crypto/sha1"
	"encoding/binary"
	"log"
)

// GetPID maps a string to one of 2^m logical points on a virtual ring.
func GetPID(s string, m int) int {
	h := sha1.New()
	if _, err := h.Write([]byte(s)); err != nil {
		log.Fatal(err)
	}
	b := h.Sum(nil)

	// Truncate down to m bits.
	pid := binary.BigEndian.Uint64(b) >> (64 - m)

	return int(pid)
}
