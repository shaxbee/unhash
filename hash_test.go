package unhash

import "testing"

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

			t.Log("hash:", v)

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
