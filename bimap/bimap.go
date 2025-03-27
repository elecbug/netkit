package bimap

type BiMap[T1 comparable, T2 comparable] struct {
	forward map[T1]T2
	reverse map[T2]T1
}

func New[T1 comparable, T2 comparable]() *BiMap[T1, T2] {
	return &BiMap[T1, T2]{
		forward: make(map[T1]T2),
		reverse: make(map[T2]T1),
	}
}

func (b *BiMap[T1, T2]) Set(key T1, value T2) {
	if oldVal, ok := b.forward[key]; ok {
		delete(b.reverse, oldVal)
	}
	if oldKey, ok := b.reverse[value]; ok {
		delete(b.forward, oldKey)
	}

	b.forward[key] = value
	b.reverse[value] = key
}

func (b *BiMap[T1, T2]) GetByKey(key T1) (T2, bool) {
	val, ok := b.forward[key]
	return val, ok
}

func (b *BiMap[T1, T2]) GetByValue(value T2) (T1, bool) {
	key, ok := b.reverse[value]
	return key, ok
}

func (b *BiMap[T1, T2]) DeleteByKey(key T1) bool {
	if val, ok := b.forward[key]; ok {
		delete(b.forward, key)
		delete(b.reverse, val)

		return true
	} else {
		return false
	}
}

func (b *BiMap[T1, T2]) DeleteByValue(value T2) bool {
	if key, ok := b.reverse[value]; ok {
		delete(b.reverse, value)
		delete(b.forward, key)

		return true
	} else {
		return false
	}
}
