package unhash

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
	}{
		{
			name: "primitive types",
			data: map[string]any{
				"bool":    true,
				"int64":   int64(1),
				"float64": float64(1.0),
				"string":  "hello",
			},
		},
		{
			name: "slices",
			data: map[string]any{
				"[]string": []string{"foo", "bar"},
				"[]any":    []any{true, int64(1), float64(1.0), "hello"},
			},
		},
		{
			name: "maps",
			data: map[string]any{
				"map[string]any": map[string]any{
					"bool":    true,
					"int64":   int64(42),
					"float64": float64(1.0),
					"string":  "hello",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			v, err := HashMap(test.data, Config{})
			if err != nil {
				t.Fatal("hash:", err)
			}

			expected, err := HashMap(test.data, Config{})
			switch {
			case err != nil:
				t.Fatal("hash:", err)
			case v != expected:
				t.Errorf("hash: expected %d, got %d", expected, v)
			}

			seeded, err := HashMap(test.data, Config{
				Seed: 42,
			})
			switch {
			case err != nil:
				t.Fatal("hash:", err)
			case v == seeded:
				t.Errorf("hash: seeded hash should return different value")
			}
		})
	}
}

func TestHashInvalid(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		expected error
	}{
		{
			name: "invalid map value type",
			data: map[string]any{
				"map": map[string]any{
					"int32": int32(1),
				},
			},
			expected: InvalidTypeError{
				Type: "int32",
				Path: "map.int32",
			},
		},
		{
			name: "invalid slice value type",
			data: map[string]any{
				"[]any": []any{int32(1)},
			},
			expected: InvalidTypeError{
				Type: "int32",
				Path: "[]any.0",
			},
		},
		{
			name: "depth limit",
			data: map[string]any{
				"map": map[string]any{
					"map": map[string]any{
						"int64": int64(1),
					},
				},
			},
			expected: MaxDepthError{
				Path: "map.map.int64",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, err := HashMap(test.data, Config{
				MaxDepth: 2,
			})
			if diff := cmp.Diff(test.expected, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("hash: %s", diff)
			}
		})
	}
}
