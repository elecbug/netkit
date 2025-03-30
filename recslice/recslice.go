package recslice

// `Node` represents either a value or a pointer to a subarray.
type Node[T any] struct {
	value *T           // Direct value
	sub   *Recslice[T] // Pointer to a subarray
}

// `Recslice` is a hybrid slice that splits into subarrays when over capacity.
type Recslice[T any] struct {
	items  []Node[T] // Slice of nodes
	maxLen int       // Maximum number of items before splitting
	length int       // Logical number of values stored
}

// `newNode` creates a Node with a value.
func newNode[T any](val T) Node[T] {
	return Node[T]{value: &val}
}

// New creates a new `Recslice` with a given `maxLen`.
func New[T any](maxLen int) *Recslice[T] {
	return &Recslice[T]{items: make([]Node[T], 0, maxLen), maxLen: maxLen}
}

// `Get` returns the value at the `index` across the entire recursive structure.
func (a *Recslice[T]) Get(index int) T {
	if index < 0 || index >= a.length {
		panic("get index out of bounds")
	}

	count := 0
	for _, item := range a.items {
		if item.value != nil {
			if count == index {
				return *item.value
			}

			count++
		} else if item.sub != nil {
			subSlice := item.sub.ToSlice() // Flattened view of sub

			if index < count+len(subSlice) {
				return subSlice[index-count]
			}

			count += len(subSlice)
		}
	}

	panic("index not found (corrupted state)")
}

// `Insert` inserts a `value` at `index`
func (a *Recslice[T]) Insert(index int, value T) {
	if index < 0 || index > a.length {
		panic("insert index out of bounds")
	}

	a.length++
	count := 0

	for i := range a.items {
		item := &a.items[i]

		if item.value != nil {
			if count == index {
				sub := New[T](a.maxLen)
				sub.subInsert(0, value)
				sub.subInsert(1, *item.value)

				a.items[i] = Node[T]{sub: sub}

				return
			}

			count++
		} else if item.sub != nil {
			subLen := item.sub.length

			if index < count+subLen {
				item.sub.Insert(index-count, value)

				return
			}

			count += subLen
		}
	}
	// If appending at the end
	a.items = append(a.items, newNode(value))
}

// `subInsert` adds a `value` at `index`, recursively creating subarrays if needed.
func (a *Recslice[T]) subInsert(index int, value T) {
	if index > len(a.items) {
		panic("insert index out of bounds")
	}

	// Insert directly if under capacity
	if len(a.items) < a.maxLen {
		a.items = append(a.items, Node[T]{})

		copy(a.items[index+1:], a.items[index:])

		a.items[index] = newNode(value)
		a.length++

		return
	}

	// If over maxLen, recursively insert into subarray
	target := a.items[index]

	// If target is a value, convert it into a new subarray
	if target.sub == nil {
		sub := New[T](a.maxLen)

		if target.value != nil {
			sub.subInsert(0, *target.value)
		}

		sub.subInsert(0, value)
		a.items[index] = Node[T]{sub: sub}

		return
	}

	// If target is already a subarray, recurse into it
	target.sub.subInsert(0, value)
}

// `Delete` deletes a `value` at `index`.
func (a *Recslice[T]) Delete(index int) {
	if index < 0 || index >= a.length {
		panic("delete index out of bounds")
	}

	a.length--
	count := 0

	for i := range a.items {
		item := &a.items[i]

		if item.value != nil {
			if count == index {
				copy(a.items[i:], a.items[i+1:])

				a.items[len(a.items)-1] = Node[T]{}
				a.items = a.items[:len(a.items)-1]

				return
			}

			count++
		} else if item.sub != nil {
			subLen := item.sub.length

			if index < count+subLen {
				item.sub.Delete(index - count)

				// Clean up empty subarray
				if item.sub.length == 0 {
					a.items[i] = Node[T]{}
				}

				return
			}
			count += subLen
		}
	}
}

// `ToSlice` flattens the entire recursive structure into a basic slice.
func (a *Recslice[T]) ToSlice() []T {
	result := make([]T, 0, a.length)
	for _, item := range a.items {
		if item.value != nil {
			result = append(result, *item.value)
		} else if item.sub != nil {
			subResult := item.sub.ToSlice()
			result = append(result, subResult...)
		}
	}
	return result
}

// Return full length of recslice
func (a Recslice[T]) Length() int {
	return a.length
}

// `Print` displays the full recursive structure of the array.
// The `toString` is formatting function.
func (a *Recslice[T]) Print(tab int, toString func(value T) string) string {
	result := "["

	for i, item := range a.items {
		if item.value != nil {
			result += toString(*item.value)
		} else if item.sub != nil {
			result += item.sub.Print(tab+1, toString)
		} else {
			result += "<nil>"
		}

		if i != len(a.items)-1 {
			result += ", "
		}
	}

	result += "]"

	return result
}
