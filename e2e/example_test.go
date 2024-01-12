package e2e

import (
	"fmt"
	"io/fs"

	"github.com/shaxbee/unhash"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

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
