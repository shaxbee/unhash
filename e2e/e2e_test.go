package e2e

import (
	"embed"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/fs"
	"path"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/shaxbee/unhash"
)

//go:embed testdata/*.yaml
var testdata embed.FS

func TestUnstructuredHash(t *testing.T) {
	objs := readObjects(t)

	for _, obj := range objs {
		obj := obj
		t.Run(obj.name, func(t *testing.T) {
			v1, err := unhash.HashMap(obj.data, unhash.Config{})
			if err != nil {
				t.Fatal(err)
			}

			v2, err := unhash.HashMap(obj.data, unhash.Config{})
			if err != nil {
				t.Fatal(err)
			}

			if v1 != v2 {
				t.Errorf("hash: expected %d, got %d", v1, v2)
			}
		})
	}
}

func ExampleHashMap() {
	data, err := fs.ReadFile(testdata, "testdata/deployment.yaml")
	if err != nil {
		panic(err)
	}

	obj := &unstructured.Unstructured{}
	if err := yaml.UnmarshalStrict(data, obj); err != nil {
		panic(err)
	}

	spec, ok, err := unstructured.NestedMap(obj.Object, "spec")
	switch {
	case err != nil:
		panic(err)
	case !ok:
		panic("no spec")
	}

	hash, err := unhash.HashMap(spec, unhash.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("hash: %x\n", hash)
	// Output: hash: b698aaf1ea3e0eec
}

func BenchmarkHash(b *testing.B) {
	objs := readObjects(b)

	b.ReportAllocs()
	b.ResetTimer()

	for _, obj := range objs {
		obj := obj
		b.Run(path.Join(obj.name, "algo=json"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data, err := json.Marshal(obj.data)
				if err != nil {
					b.Fatal("marshal:", err)
				}

				hash := fnv.New64()
				_, _ = hash.Write(data)
			}
		})
		b.Run(path.Join(obj.name, "algo=unstructured"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := unhash.HashMap(obj.data, unhash.Config{})
				if err != nil {
					b.Fatal(i, err)
				}
			}
		})
	}
}

type object struct {
	name string
	data map[string]any
}

func readObjects(t testing.TB) []object {
	fsys, err := fs.Sub(testdata, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	files, err := fs.Glob(fsys, "*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	res := make([]object, len(files))
	for i, filename := range files {
		data, err := fs.ReadFile(fsys, filename)
		if err != nil {
			t.Fatal(err)
		}

		obj := &unstructured.Unstructured{}
		if err := yaml.UnmarshalStrict(data, obj); err != nil {
			t.Fatalf("unmarshal %q: %v", filename, err)
		}

		res[i] = object{
			name: filename,
			data: obj.Object,
		}
	}

	return res
}
