package redblacktree

// Iterator holding the iterator`s state
type Iterator struct {
	t *Tree
	n *Node
	p position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// Iterator returns a stateful iterator key/value pairs.
func (tr *Tree) Iterator() Iterator {
	return Iterator{t: tr, n: nil, p: begin}
}

// Next moves the iterator to the next element
func (it *Iterator) Next() bool {
	if it.p == end {
		goto end
	}

	if it.p == begin {
		l := it.t.left()
		if l == nil {
			goto end
		}
		it.n = l
		goto between
	}
	if it.n.right != nil {
		it.n = it.n.right
		for it.n.left != nil {
			it.n = it.n.left
		}
		goto between
	}
	if it.n.parent != nil {
		n := it.n
		for it.n.parent != nil {
			it.n = it.n.parent
			if it.t.compare(n.key, it.n.key) >= 0 {
				goto between
			}
		}
	}
end:
	it.n = nil
	it.p = end
	return false
between:
	it.p = between
	return true
}

// Prev moves the iterator to the prev element
func (it *Iterator) Prev() bool {
	if it.p == begin {
		goto begin
	}
	if it.p == end {
		r := it.t.right()
		if r == nil {
			goto begin
		}
		it.n = r
		goto between
	}
	if it.n.left != nil {
		it.n = it.n.left
		for it.n.right != nil {
			it.n = it.n.right
		}
		goto between
	}
	if it.n.parent != nil {
		n := it.n
		for it.n.parent != nil {
			it.n = it.n.parent
			if it.t.compare(n.key, it.n.key) >= 0 {
				goto between
			}
		}
	}
begin:
	it.n = nil
	it.p = begin
	return false

between:
	it.p = between
	return true
}

// It returns current Node
func (it *Iterator) It() *Node {
	return it.n
}

// Value returns the current element`s value.
func (it *Iterator) Value() interface{} {
	return it.n.val
}

// Key returns the current element`s key.
func (it *Iterator) Key() interface{} {
	return it.n.key
}

// Begin resets the iterator to its start(one-before-first)
func (it *Iterator) Begin() {
	it.n = nil
	it.p = begin
}

// End resets the iterator to its end (one-past-the-end)
func (it *Iterator) End() {
	it.n = nil
	it.p = end
}

// First moves the iterator to the first element
func (it *Iterator) First() bool {
	it.Begin()
	return it.Next()
}

// Last moves the iterator to the last element
func (it *Iterator) Last() bool {
	it.End()
	return it.Prev()
}
