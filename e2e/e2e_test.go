package e2e

import (
	"embed"
	"io/fs"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/shaxbee/unhash"
)

//go:embed testdata/*.yaml
var testdata embed.FS

func TestUnstructuredHash(t *testing.T) {
	files := loadTestData(t, testdata, "testdata/*.yaml")

	for filename, obj := range files {
		filename, obj := filename, obj
		t.Run(filename, func(t *testing.T) {
			v1, err := unhash.HashMap(obj.Object, unhash.Config{})
			if err != nil {
				t.Fatal(err)
			}

			v2, err := unhash.HashMap(obj.Object, unhash.Config{})
			if err != nil {
				t.Fatal(err)
			}

			if v1 != v2 {
				t.Errorf("hash: expected %d, got %d", v1, v2)
			}
		})
	}
}

func BenchmarkUnstructuredHash(b *testing.B) {
	files := loadTestData(b, testdata, "testdata/*.yaml")

	b.ReportAllocs()
	b.ResetTimer()

	for filename, obj := range files {
		filename, obj := filename, obj

		b.Run(filename, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := unhash.HashMap(obj.Object, unhash.Config{})
				if err != nil {
					b.Fatal(i, err)
				}
			}
		})
	}
}

func loadTestData(t testing.TB, fsys fs.FS, pattern string) map[string]*unstructured.Unstructured {
	files, err := fs.Glob(fsys, pattern)
	if err != nil {
		t.Fatal(err)
	}

	res := make(map[string]*unstructured.Unstructured, len(files))
	for _, filename := range files {
		data, err := fs.ReadFile(fsys, filename)
		if err != nil {
			t.Fatal(err)
		}

		obj := &unstructured.Unstructured{}
		if err := yaml.UnmarshalStrict(data, obj); err != nil {
			t.Fatalf("unmarshal %q: %v", filename, err)
		}

		res[filename] = obj
	}

	return res
}
