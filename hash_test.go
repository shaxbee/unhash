package unhash

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		mutation func(map[string]any)
	}{
		{
			name: "primitive types",
			data: map[string]any{
				"bool":    true,
				"int64":   int64(1),
				"float64": float64(1.0),
				"string":  "hello",
			},
			mutation: func(data map[string]any) {
				data["bool"] = false
			},
		},
		{
			name: "slices",
			data: map[string]any{
				"[]any": []any{true, int64(1), float64(1.0), "hello", nil},
			},
			mutation: func(data map[string]any) {
				data["[]any"] = append(data["[]any"].([]any), "world")
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
			mutation: func(data map[string]any) {
				data = data["map[string]any"].(map[string]any)
				data["bool"] = false
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			expected, err := HashMap(test.data, Config{})
			if err != nil {
				t.Fatal("hash:", err)
			}

			hash, err := HashMap(test.data, Config{})
			switch {
			case err != nil:
				t.Fatal("hash:", err)
			case hash != expected:
				t.Errorf("hash: expected %d, got %d", expected, hash)
			}

			hash, err = HashMap(test.data, Config{
				Seed: 42,
			})
			switch {
			case err != nil:
				t.Fatal("hash:", err)
			case hash == expected:
				t.Errorf("hash: seeded hash should return different value")
			}

			if test.mutation != nil {
				test.mutation(test.data)
			}

			hash, err = HashMap(test.data, Config{})
			switch {
			case err != nil:
				t.Fatal("hash:", err)
			case hash == expected:
				t.Errorf("hash: mutated data should return different value")
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
				Path: "map/int32",
			},
		},
		{
			name: "invalid slice value type",
			data: map[string]any{
				"[]any": []any{int32(1)},
			},
			expected: InvalidTypeError{
				Type: "int32",
				Path: "[]any/0",
			},
		},
		{
			name: "nested map depth limit",
			data: map[string]any{
				"map": map[string]any{
					"map": map[string]any{
						"int64": int64(1),
					},
				},
			},
			expected: MaxDepthError{
				Path: "map/map/int64",
			},
		},
		{
			name: "nested slice depth limit",
			data: map[string]any{
				"map": map[string]any{
					"slice": []any{int64(1), float64(1.0)},
				},
			},
			expected: MaxDepthError{
				Path: "map/slice/0",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, err := HashMap(test.data, Config{
				MaxDepth: 2,
			})
			if err == nil || err.Error() != test.expected.Error() {
				t.Errorf("hash: expected %v, got %v", test.expected, err)
			}
		})
	}
}
