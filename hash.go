package unhash

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"unsafe"
)

var ErrMaxDepth = fmt.Errorf("unhash: max depth reached")

func HashMap(data map[string]any, config Config) (uint64, error) {
	config = ConfigDefault(config)

	res, err := visitMap(data, config)
	if err != nil {
		return 0, err
	}

	h := config.Hash()
	if config.Seed != 0 {
		writeUint64(h, config.Seed)
	}
	writeUint64(h, res)

	return h.Sum64(), nil
}

func visitMap(data map[string]any, config Config) (uint64, error) {
	h := config.Hash()

	var sum uint64
	for k, v := range data {
		h.Reset()

		res, err := visitValue(v, config)
		if err != nil {
			return 0, err
		}

		writeString(h, k)
		writeUint64(h, res)

		sum ^= h.Sum64()
	}

	h.Reset()
	writeUint64(h, sum)

	return h.Sum64(), nil
}

func visitSlice(data []any, config Config) (uint64, error) {
	h := config.Hash()

	for _, v := range data {
		res, err := visitValue(v, config)
		if err != nil {
			return 0, err
		}

		writeUint64(h, res)
	}

	return h.Sum64(), nil
}

func visitValue(data any, config Config) (uint64, error) {
	h := config.Hash()

	switch tv := data.(type) {
	case string:
		writeString(h, tv)
	case int64:
		writeUint64(h, uint64(tv))
	case float64:
		writeUint64(h, math.Float64bits(tv))
	case bool:
		var bv uint64
		if tv {
			bv = 1
		}
		writeUint64(h, bv)
	case []string:
		for _, s := range tv {
			writeString(h, s)
		}
	case []any:
		res, err := visitSlice(tv, config)
		if err != nil {
			return 0, err
		}
		writeUint64(h, res)
	case map[string]any:
		res, err := visitMap(tv, config)
		if err != nil {
			return 0, err
		}
		writeUint64(h, res)
	default:
		return 0, fmt.Errorf("unhash: unsupported type %T", tv)
	}

	return h.Sum64(), nil
}

func writeString(w io.Writer, data string) {
	w.Write(unsafe.Slice(unsafe.StringData(data), len(data)))
}

func writeUint64(w io.Writer, v uint64) {
	b := [8]byte{}
	binary.NativeEndian.PutUint64(b[:], v)
	w.Write(b[:])
}
