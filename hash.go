package unhash

import (
	"fmt"
	"hash/fnv"

	"github.com/shaxbee/unhash/internal/fasthash/fnv1"
)

var (
	DefaultHash     = fnv.New64
	DefaultMaxDepth = 20
)

type Config struct {
	MaxDepth int
	Seed     uint64
}

// MaxDepthError indicates value that is nested above max depth
type MaxDepthError struct {
	Path string
}

func (e MaxDepthError) Error() string {
	return fmt.Sprintf("unhash: max depth reached at %s", e.Path)
}

// InvalidValueTypeError
type InvalidTypeError struct {
	Path string
	Type string
}

func (e InvalidTypeError) Error() string {
	return fmt.Sprintf("unhash: unsupported type %s", e.Type)
}

// HashMap computes a hash of unstructured object.
//
// Supported values:
//   - string
//   - int64
//   - float64
//   - bool
//   - map[string]any
//   - []any
func HashMap(data map[string]any, config Config) (uint64, error) {
	if config.MaxDepth == 0 {
		config.MaxDepth = DefaultMaxDepth
	}

	v := visitor{
		config: config,
	}

	res, err := v.visitMap(data)
	if err != nil {
		return 0, err
	}

	var hash = fnv1.Init
	if config.Seed != 0 {
		hash = fnv1.AddUint64(hash, config.Seed)
	}
	hash = fnv1.AddUint64(hash, res)

	return hash, nil
}
