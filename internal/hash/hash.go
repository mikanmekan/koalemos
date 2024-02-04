package hash

import (
	"hash/fnv"
)

func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s)) // Does not return err.
	return h.Sum64()
}
