// Package bimap implements a bidirectional map with constant-time lookups
// from key to value and value to key.
package bimap

// Bimap is a bidirectional map providing O(1) lookup for both directions.
type Bimap[T1 comparable, T2 comparable] struct {
	forward map[T1]T2
	reverse map[T2]T1
}

// Pair holds a single key/value pair.
type Pair[T1 comparable, T2 comparable] struct {
	Key   T1
	Value T2
}

// New creates an empty Bimap.
func New[T1 comparable, T2 comparable]() *Bimap[T1, T2] {
	return &Bimap[T1, T2]{
		forward: make(map[T1]T2),
		reverse: make(map[T2]T1),
	}
}

// Set inserts or updates a key/value pair.
// If the key or the value already exists, the previous association is removed
// so that keys and values remain unique.
func (b *Bimap[T1, T2]) Set(key T1, value T2) {
	if oldVal, ok := b.forward[key]; ok {
		delete(b.reverse, oldVal)
	}
	if oldKey, ok := b.reverse[value]; ok {
		delete(b.forward, oldKey)
	}

	b.forward[key] = value
	b.reverse[value] = key
}

// GetByKey returns the value for key. The boolean is false if key is not found.
func (b *Bimap[T1, T2]) GetByKey(key T1) (T2, bool) {
	val, ok := b.forward[key]
	return val, ok
}

// GetByValue returns the key for value. The boolean is false if value is not found.
func (b *Bimap[T1, T2]) GetByValue(value T2) (T1, bool) {
	key, ok := b.reverse[value]
	return key, ok
}

// DeleteByKey removes the pair identified by key.
// It returns true if the pair existed and was removed.
func (b *Bimap[T1, T2]) DeleteByKey(key T1) bool {
	if val, ok := b.forward[key]; ok {
		delete(b.forward, key)
		delete(b.reverse, val)

		return true
	} else {
		return false
	}
}

// DeleteByValue removes the pair identified by value.
// It returns true if the pair existed and was removed.
func (b *Bimap[T1, T2]) DeleteByValue(value T2) bool {
	if key, ok := b.reverse[value]; ok {
		delete(b.reverse, value)
		delete(b.forward, key)

		return true
	} else {
		return false
	}
}

// ToList returns all pairs in the map as a slice.
func (b Bimap[T1, T2]) ToList() []Pair[T1, T2] {
	result := make([]Pair[T1, T2], len(b.forward))
	i := 0

	for k, v := range b.forward {
		result[i] = Pair[T1, T2]{
			Key:   k,
			Value: v,
		}

		i++
	}

	return result
}
