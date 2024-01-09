package unhash

import (
	"hash"
	"hash/fnv"
)

var (
	DefaultHash     = fnv.New64
	DefaultMaxDepth = 20
)

type Config struct {
	Hash     func() hash.Hash64
	MaxDepth int
	Seed     uint64
}

func ConfigDefault(c Config) Config {
	if c.Hash == nil {
		c.Hash = DefaultHash
	}

	if c.MaxDepth == 0 {
		c.MaxDepth = DefaultMaxDepth
	}

	return c
}
