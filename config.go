package unhash

import (
	"hash/fnv"
)

var (
	DefaultHash     = fnv.New64
	DefaultMaxDepth = 20
)

type Config struct {
	MaxDepth int
	Seed     uint64
}

func ConfigDefault(c Config) Config {
	if c.MaxDepth == 0 {
		c.MaxDepth = DefaultMaxDepth
	}

	return c
}
