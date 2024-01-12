# Unhash

[![Go Reference](https://pkg.go.dev/badge/github.com/shaxbee/unhash.svg)](https://pkg.go.dev/github.com/shaxbee/unhash)
[![Go Coverage](https://github.com/shaxbee/unhash/wiki/coverage.svg)](https://raw.githack.com/wiki/shaxbee/unhash/coverage.html)

Unhash is a hash function for [unstructured](k8s.io/apimachinery/pkg/apis/meta/v1/unstructured) data.

It is designed to avoid allocations and utilizes FNV1 hash implementation from [fasthash](https://pkg.go.dev/github.com/segmentio/fasthash/fnv1).


## Usage

### Supported values

- `string`
- `int64`
- `float64`
- `bool`
- `map[string]any`
- `[]any`

### Example

```go
data, err := fs.ReadFile(testdata, "testdata/deployment.yaml")
if err != nil {
    panic(err)
}

obj := &unstructured.Unstructured{}
if err := yaml.UnmarshalStrict(data, obj); err != nil {
    panic(err)
}

// extract deployment spec
spec, ok, err := unstructured.NestedMap(obj.Object, "spec")
switch {
case err != nil:
    panic(err)
case !ok:
    panic("no spec")
}

// hash spec
hash, err := unhash.HashMap(spec, unhash.Config{})
if err != nil {
    panic(err)
}
```

## Benchmark

Performance compared to json encode + fnv1 hash:

```
                        │     json     │            unstructured             │
                        │    sec/op    │   sec/op     vs base                │
Hash/deployment.yaml-10   10.190µ ± 1%   1.946µ ± 2%  -80.90% (p=0.000 n=10)
Hash/issuer.yaml-10       1347.0n ± 1%   269.1n ± 2%  -80.02% (p=0.000 n=10)
Hash/pod.yaml-10          56.642µ ± 1%   9.964µ ± 1%  -82.41% (p=0.000 n=10)
geomean                    9.195µ        1.735µ       -81.14%

                        │     json     │            unstructured            │
                        │     B/op     │    B/op     vs base                │
Hash/deployment.yaml-10    7243.0 ± 0%   744.0 ± 0%  -89.73% (p=0.000 n=10)
Hash/issuer.yaml-10        912.00 ± 0%   72.00 ± 0%  -92.11% (p=0.000 n=10)
Hash/pod.yaml-10          36477.5 ± 0%   744.0 ± 0%  -97.96% (p=0.000 n=10)
geomean                   6.077Ki        341.6       -94.51%

                        │     json     │            unstructured            │
                        │  allocs/op   │ allocs/op   vs base                │
Hash/deployment.yaml-10   157.000 ± 0%   5.000 ± 0%  -96.82% (p=0.000 n=10)
Hash/issuer.yaml-10        22.000 ± 0%   2.000 ± 0%  -90.91% (p=0.000 n=10)
Hash/pod.yaml-10          673.000 ± 0%   5.000 ± 0%  -99.26% (p=0.000 n=10)
geomean                     132.5        3.684       -97.22%
```