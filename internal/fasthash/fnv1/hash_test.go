// Implementation from https://github.com/segmentio/fasthash

package fnv1

import (
	"encoding/binary"
	"hash/fnv"
	"testing"
)

func TestAddString(t *testing.T) {
	for _, s := range referenceStrings {
		reference := fnv.New64()
		_, _ = reference.Write([]byte(s))
		expected := reference.Sum64()

		sum := AddString(Init, s)
		if sum != expected {
			t.Errorf("invalid hash for %q: expected %x but got %x", s, expected, sum)
		}
	}
}

func TestAddBytes(t *testing.T) {
	for _, s := range referenceStrings {
		reference := fnv.New64()
		_, _ = reference.Write([]byte(s))
		expected := reference.Sum64()

		sum := AddBytes(Init, []byte(s))
		if sum != expected {
			t.Errorf("invalid hash for %q: expected %x but got %x", s, expected, sum)
		}
	}
}

func TestAddUint64(t *testing.T) {
	reference := fnv.New64()
	b := [8]byte{}
	binary.BigEndian.PutUint64(b[:], 42)
	_, _ = reference.Write(b[:])
	expected := reference.Sum64()

	sum := AddUint64(Init, 42)
	if sum != expected {
		t.Errorf("invalid hash: expected %x but got %x", expected, sum)
	}
}

var referenceStrings = []string{
	"",
	"A",
	"hey",
	"Hello World!",
	"DAB45194-42CC-4106-AB9F-2447FA4D35C2",
	"你好吗",
}
