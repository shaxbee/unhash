package unhash

import (
	"fmt"

	"github.com/shaxbee/unhash/internal/fasthash/fnv1"
)

type MaxDepthError struct {
	Path string
}

func (e MaxDepthError) Error() string {
	return fmt.Sprintf("unhash: max depth reached at %s", e.Path)
}

type InvalidTypeError struct {
	Path string
	Type string
}

func (e InvalidTypeError) Error() string {
	return fmt.Sprintf("unhash: unsupported type %s", e.Type)
}

func HashMap(data map[string]any, config Config) (uint64, error) {
	config = ConfigDefault(config)

	v := visitor{
		config: config,
	}

	res, err := v.visitMap(data)
	if err != nil {
		return 0, err
	}

	var hash = fnv1.Init
	if config.Seed != 0 {
		fnv1.AddUint64(hash, config.Seed)
	}
	fnv1.AddUint64(hash, res)

	return hash, nil
}
