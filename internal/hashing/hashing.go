package hashing

import (
    "crypto/sha1"
    "encoding/binary"
)

func MHash(s string, m int) int {
    h := sha1.New()
    b := h.Sum([]byte(s))

    // Truncate down to m bits.
    pid := binary.BigEndian.Uint64(b) & ((1 << m) - 1)

    return int(pid)
}

func BHash(b []byte) []byte {
    h := sha1.New()
    return h.Sum([]byte(b))
}
