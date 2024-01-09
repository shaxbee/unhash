package unhash

import (
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/shaxbee/unhash/internal/fasthash/fnv1"
)

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

		var hash = fnv1.Init
		hash = fnv1.AddString(hash, key)
		hash = fnv1.AddUint64(hash, res)

		sum ^= hash

		v.pop()
	}

	var hash = fnv1.Init
	return fnv1.AddUint64(hash, sum), nil
}

func (v *visitor) visitSlice(data []any) (uint64, error) {
	var hash = fnv1.Init
	for idx, value := range data {
		if err := v.push(segment{idx: idx}); err != nil {
			return 0, err
		}

		res, err := v.visitValue(value)
		if err != nil {
			return 0, err
		}

		hash = fnv1.AddUint64(hash, res)

		v.pop()
	}

	return hash, nil
}

func (v *visitor) visitValue(data any) (uint64, error) {
	if data == nil {
		return fnv1.Init, nil
	}

	var hash = fnv1.Init
	switch tv := data.(type) {
	case string:
		hash = fnv1.AddString(hash, tv)
	case int64:
		hash = fnv1.AddUint64(hash, uint64(tv))
	case float64:
		hash = fnv1.AddUint64(hash, math.Float64bits(tv))
	case bool:
		var bv uint64
		if tv {
			bv = 1
		}
		hash = fnv1.AddUint64(hash, bv)
	case []any:
		res, err := v.visitSlice(tv)
		if err != nil {
			return 0, err
		}
		hash = fnv1.AddUint64(hash, res)
	case map[string]any:
		res, err := v.visitMap(tv)
		if err != nil {
			return 0, err
		}
		hash = fnv1.AddUint64(hash, res)
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
	v.path = v.path[:len(v.path)-1]
}

func (v *visitor) push(seg segment) error {
	v.path = append(v.path, seg)
	if v.config.MaxDepth > 0 && len(v.path) > v.config.MaxDepth {
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
