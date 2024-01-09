package unhash

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/segmentio/fasthash/fnv1"
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

	var hash = config.Seed
	fnv1.AddUint64(hash, res)

	return hash, nil
}

type visitor struct {
	path   []segment
	config Config
}

type segment struct {
	str string
	idx int
}

func (v *visitor) visitMap(data map[string]any) (uint64, error) {
	var sum uint64
	for key, value := range data {
		if err := v.push(segment{str: key}); err != nil {
			return 0, err
		}

		res, err := v.visitValue(value)
		if err != nil {
			return 0, err
		}

		var hash uint64
		fnv1.AddString64(hash, key)
		fnv1.AddUint64(hash, res)

		sum ^= hash

		v.pop()
	}

	return fnv1.HashUint64(sum), nil
}

func (v *visitor) visitSlice(data []any) (uint64, error) {
	var hash uint64
	for idx, value := range data {
		if err := v.push(segment{idx: idx}); err != nil {
			return 0, err
		}

		res, err := v.visitValue(value)
		if err != nil {
			return 0, err
		}

		fnv1.AddUint64(hash, res)

		v.pop()
	}

	return hash, nil
}

func (v *visitor) visitValue(data any) (uint64, error) {
	if data == nil {
		return uint64(0), nil
	}

	var hash uint64
	switch tv := data.(type) {
	case string:
		fnv1.AddString64(hash, tv)
	case int64:
		fnv1.AddUint64(hash, uint64(tv))
	case float64:
		fnv1.AddUint64(hash, math.Float64bits(tv))
	case bool:
		var bv uint64
		if tv {
			bv = 1
		}
		fnv1.AddUint64(hash, bv)
	case []string:
		for _, s := range tv {
			fnv1.AddString64(hash, s)
		}
	case []any:
		res, err := v.visitSlice(tv)
		if err != nil {
			return 0, err
		}
		fnv1.AddUint64(hash, res)
	case map[string]any:
		res, err := v.visitMap(tv)
		if err != nil {
			return 0, err
		}
		fnv1.AddUint64(hash, res)
	default:
		tpe := reflect.TypeOf(data).String()
		return 0, InvalidTypeError{
			Path: v.current(),
			Type: tpe,
		}
	}

	return hash, nil
}

func (v *visitor) current() string {
	var path []string
	for _, seg := range v.path {
		path = append(path, seg.String())
	}

	return strings.Join(path, ".")
}

func (v *visitor) pop() {
	if len(v.path) == 0 {
		panic("unhash: empty path")
	}

	v.path = v.path[:len(v.path)-1]
}

func (v *visitor) push(seg segment) error {
	v.path = append(v.path, seg)
	if len(v.path) > v.config.MaxDepth {
		return MaxDepthError{
			Path: v.current(),
		}
	}

	return nil
}

func (s segment) String() string {
	if s.str != "" {
		return s.str
	}

	return strconv.Itoa(s.idx)
}
