// Implementation from https://github.com/segmentio/fasthash

package fnv1

const (
	// FNV-1
	offset64 = uint64(14695981039346656037)
	prime64  = uint64(1099511628211)

	// Init is what 64 bits hash values should be initialized with.
	Init = offset64
)

// AddString adds the hash of s to the precomputed hash value h.
func AddString(h uint64, s string) uint64 {
	for len(s) >= 8 {
		h = (h * prime64) ^ uint64(s[0])
		h = (h * prime64) ^ uint64(s[1])
		h = (h * prime64) ^ uint64(s[2])
		h = (h * prime64) ^ uint64(s[3])
		h = (h * prime64) ^ uint64(s[4])
		h = (h * prime64) ^ uint64(s[5])
		h = (h * prime64) ^ uint64(s[6])
		h = (h * prime64) ^ uint64(s[7])
		s = s[8:]
	}

	if len(s) >= 4 {
		h = (h * prime64) ^ uint64(s[0])
		h = (h * prime64) ^ uint64(s[1])
		h = (h * prime64) ^ uint64(s[2])
		h = (h * prime64) ^ uint64(s[3])
		s = s[4:]
	}

	if len(s) >= 2 {
		h = (h * prime64) ^ uint64(s[0])
		h = (h * prime64) ^ uint64(s[1])
		s = s[2:]
	}

	if len(s) > 0 {
		h = (h * prime64) ^ uint64(s[0])
	}

	return h
}

// AddBytes adds the hash of b to the precomputed hash value h.
func AddBytes(h uint64, b []byte) uint64 {
	for len(b) >= 8 {
		h = (h * prime64) ^ uint64(b[0])
		h = (h * prime64) ^ uint64(b[1])
		h = (h * prime64) ^ uint64(b[2])
		h = (h * prime64) ^ uint64(b[3])
		h = (h * prime64) ^ uint64(b[4])
		h = (h * prime64) ^ uint64(b[5])
		h = (h * prime64) ^ uint64(b[6])
		h = (h * prime64) ^ uint64(b[7])
		b = b[8:]
	}

	if len(b) >= 4 {
		h = (h * prime64) ^ uint64(b[0])
		h = (h * prime64) ^ uint64(b[1])
		h = (h * prime64) ^ uint64(b[2])
		h = (h * prime64) ^ uint64(b[3])
		b = b[4:]
	}

	if len(b) >= 2 {
		h = (h * prime64) ^ uint64(b[0])
		h = (h * prime64) ^ uint64(b[1])
		b = b[2:]
	}

	if len(b) > 0 {
		h = (h * prime64) ^ uint64(b[0])
	}

	return h
}

// AddUint64 adds the hash value of the 8 bytes of u to h.
func AddUint64(h uint64, u uint64) uint64 {
	h = (h * prime64) ^ ((u >> 56) & 0xFF)
	h = (h * prime64) ^ ((u >> 48) & 0xFF)
	h = (h * prime64) ^ ((u >> 40) & 0xFF)
	h = (h * prime64) ^ ((u >> 32) & 0xFF)
	h = (h * prime64) ^ ((u >> 24) & 0xFF)
	h = (h * prime64) ^ ((u >> 16) & 0xFF)
	h = (h * prime64) ^ ((u >> 8) & 0xFF)
	h = (h * prime64) ^ ((u >> 0) & 0xFF)
	return h
}
