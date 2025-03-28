package bimap

// Bidirectional map structure to ensure constant time complexity.
type Bimap[T1 comparable, T2 comparable] struct {
	forward map[T1]T2
	reverse map[T2]T1
}

// `Pair` structure consisting of one `Key`, `Value` pair
type Pair[T1 comparable, T2 comparable] struct {
	Key   T1
	Value T2
}

// Create a `Bimap` with `T1` and `T2` as types.
func New[T1 comparable, T2 comparable]() *Bimap[T1, T2] {
	return &Bimap[T1, T2]{
		forward: make(map[T1]T2),
		reverse: make(map[T2]T1),
	}
}

// Save data to `Bimap`.
// It allows overwriting, and removes the existing pair if any of the
// `key` and `value` that were previously used overlap.
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

// Get the value from the `key`.
// if it exists `(value, true)` and if it doesn't, it's `false`.
func (b *Bimap[T1, T2]) GetByKey(key T1) (T2, bool) {
	val, ok := b.forward[key]
	return val, ok
}

// Get the key from the `value`.
// if it exists `(key, true)` and if it doesn't, it's `false`.
func (b *Bimap[T1, T2]) GetByValue(value T2) (T1, bool) {
	key, ok := b.reverse[value]
	return key, ok
}

// Remove the pair of bidirectional as the `key`.
// if it remove `true` and if doesn't, it's `false`.
func (b *Bimap[T1, T2]) DeleteByKey(key T1) bool {
	if val, ok := b.forward[key]; ok {
		delete(b.forward, key)
		delete(b.reverse, val)

		return true
	} else {
		return false
	}
}

// Remove the pair of bidirectional as the `value`.
// if it remove `true` and if doesn't, it's `false`.
func (b *Bimap[T1, T2]) DeleteByValue(value T2) bool {
	if key, ok := b.reverse[value]; ok {
		delete(b.reverse, value)
		delete(b.forward, key)

		return true
	} else {
		return false
	}
}

// Returns `Bimap` as a list of simple `Key`, `Value` `Pair`.
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
