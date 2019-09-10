package redblacktree

import (
	"fmt"

	"github.com/yamakiller/magicNet/st/comparator"
)

const (
	colorRed   = true
	colorBlock = false

	alloterMalloc = true
	alloterFree   = false
)

// Node red-black Tree is Node
type Node struct {
	//node
	parent, left, right *Node
	//property
	key, val interface{}
	color    bool
}

func (n *Node) grandparent() *Node {
	if n != nil && n.parent != nil {
		return n.parent.parent
	}
	return nil
}

func (n *Node) uncle() *Node {
	if n == nil || n.parent == nil || n.parent.parent == nil {
		return nil
	}

	return n.parent.sibling()
}

func (n *Node) sibling() *Node {
	if n == nil || n.parent == nil {
		return nil
	}
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

func (n *Node) maximumNode() *Node {
	if n == nil {
		return nil
	}
	for n.right != nil {
		n = n.right
	}
	return n
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.key)
}

// Alloter Tree alloter
type Alloter struct {
	A func(k, v interface{}, c bool) *Node
	F func(p *Node)
}

// Tree red-black tree
type Tree struct {
	root    *Node
	size    int
	compare comparator.Comparator
	alloter *Alloter
}

// NewTree instantiates a red-black tree with the custom comparator
func NewTree(comparator comparator.Comparator, alloc *Alloter) *Tree {
	return &Tree{compare: comparator, alloter: alloc}
}

// NewTreeIntComparator instantiates a red-black tree with the IntComparator, i.e. keys are of type int.
func NewTreeIntComparator(alloc *Alloter) *Tree {
	return &Tree{compare: comparator.IntComparator, alloter: alloc}
}

// NewTreeStringComparator instantiates a red-black tree with the IntComparator, i.e. keys are of type string.
func NewTreeStringComparator(alloc *Alloter) *Tree {
	return &Tree{compare: comparator.StringComparator, alloter: alloc}
}

var (
	defaultAlloter = Alloter{A: func(k, v interface{}, c bool) *Node {
		return &Node{key: k, val: v, color: c}
	},
		F: func(p *Node) {

		},
	}
)

func (t *Tree) newNode(k, v interface{}, c bool) *Node {
	if t.alloter == nil {
		return defaultAlloter.A(k, v, c)
	}

	return t.alloter.A(k, v, c)
}

func (t *Tree) freeNode(n *Node) {
	if t.alloter == nil {
		defaultAlloter.F(n)
		return
	}

	t.alloter.F(n)
	return
}

func (t *Tree) isRed(n *Node) bool {
	return n != nil && n.color
}

// Empty returns true if tree does not contain any nodes
func (t *Tree) Empty() bool {
	return (t.size == 0)
}

// Size returns number of nodes in the tree
func (t *Tree) Size() int {
	return t.size
}

// Clear removes all nodes from the tree
func (t *Tree) Clear() {
	//! With memory management, nodes need to be released one by one
	if t.alloter != nil {
		it := t.Iterator()
		for i := 0; it.Next(); i++ {
			if it.It() != nil {
				t.freeNode(it.It())
			}
		}
	}

	t.root = nil
	t.size = 0
}

func (t *Tree) String() string {
	str := "RedBlackTree\n"
	if !t.Empty() {
		output(t.root, "", true, &str)
	}
	return str
}

func output(n *Node, prefix string, isTail bool, str *string) {
	if n.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(n.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += n.String() + "\n"
	if n.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(n.left, newPrefix, true, str)
	}
}

func (t *Tree) getDepth() int {
	var depthCall func(n *Node) int
	depthCall = func(n *Node) int {
		if n == nil {
			return 0
		}

		if n.left == nil && n.right == nil {
			return 1
		}

		ld := depthCall(n.left)
		rd := depthCall(n.right)

		if ld > rd {
			return ld + 1
		}
		return rd + 1
	}

	return depthCall(t.root)
}

func (t *Tree) leftRotate(n *Node) {
	if n.right == nil {
		return
	}

	r := n.right
	n.right = r.left
	if n.right != nil {
		n.right.parent = n
	}

	r.parent = n.parent
	if n.parent == nil {
		t.root = r
	} else {
		if n.parent.left == n {
			n.parent.left = r
		} else {
			n.parent.right = r
		}
	}
	r.left = n
	n.parent = r
}

func (t *Tree) rightRotate(n *Node) {
	if n.left == nil {
		return
	}

	l := n.left
	n.left = l.right
	if n.left != nil {
		n.left.parent = n
	}

	l.parent = n.parent
	if n.parent == nil {
		t.root = l
	} else {
		if n.parent.left == n {
			n.parent.left = l
		} else {
			n.parent.right = l
		}
	}
	l.right = n
	n.parent = l
}

func (t *Tree) lookup(key interface{}) *Node {
	n := t.root
	for n != nil {
		compare := t.compare(key, n.key)
		if compare > 0 {
			n = n.right
		} else if compare < 0 {
			n = n.left
		} else {
			return n
		}
	}
	return nil
}

// Left returns the left-most (min) node or nil if tree is empty
func (t *Tree) left() *Node {
	var parent *Node
	cur := t.root
	for cur != nil {
		parent = cur
		cur = cur.left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty
func (t *Tree) right() *Node {
	var parent *Node
	cur := t.root
	for cur != nil {
		parent = cur
		cur = cur.right
	}
	return parent
}

// Keys returns all keys
func (t *Tree) Keys() []interface{} {
	keys := make([]interface{}, t.size)
	it := t.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values
func (t *Tree) Values() []interface{} {
	vals := make([]interface{}, t.size)
	it := t.Iterator()
	for i := 0; it.Next(); i++ {
		vals[i] = it.Value()
	}
	return vals
}

// Get sreach the node in thre tree by key and returns its value or nil if key is not found in tree
func (t *Tree) Get(k interface{}) (interface{}, bool) {
	n := t.lookup(k)
	if n != nil {
		return n.val, true
	}
	return nil, false
}

// Insert inserts node into the tree
func (t *Tree) Insert(k, v interface{}) {
	var n *Node
	if t.Empty() {
		t.root = t.newNode(k, v, colorRed)
		return
	}

	cur := t.root
	for {
		if t.compare(k, cur.key) > 0 {
			if cur.right == nil {
				cur.right = t.newNode(k, v, colorRed)
				n = cur.right
				break
			}

			cur = cur.right
			continue
		} else if t.compare(k, cur.key) < 0 {
			if cur.left == nil {
				cur.left = t.newNode(k, v, colorRed)
				n = cur.left
				break
			}

			cur = cur.left
			continue
		}

		cur.key = k
		cur.val = v
		return
	}
	n.parent = cur

	t.insertFixup1(n)
	t.size++
}

// Erase remove the node from the tree by key.
func (t *Tree) Erase(k interface{}) {
	var c *Node
	n := t.lookup(k)
	if n == nil {
		return
	}
	if n.left != nil && n.right != nil {
		pred := n.left.maximumNode()
		n.key = pred.key
		n.val = pred.val
		n = pred
	}
	if n.left == nil || n.right == nil {
		if n.right == nil {
			c = n.left
		} else {
			c = n.right
		}
		if n.color == colorBlock {
			n.color = t.isRed(c)
			t.eraseFixup1(n)
		}
		if n.parent == nil {
			t.root = c
		} else {
			if n == n.left.parent {
				n.parent.left = c
			} else {
				n.parent.right = c
			}
		}
		if c != nil {
			c.parent = n.parent
		}
		if n.parent == nil && c != nil {
			c.color = colorBlock
		}
	}
	t.size--
	if n.parent == nil &&
		n.left == nil &&
		n.right == nil {
		// ?  Logging and testing. Can it be performed here?
		t.freeNode(n)
	}
}

func (t *Tree) insertFixup1(n *Node) {
	if n.parent == nil {
		n.color = colorBlock
	} else {
		t.insertFixup2(n)
	}
}

func (t *Tree) insertFixup2(n *Node) {
	if !t.isRed(n) {
		return
	}
	t.insertFixup3(n)
}

func (t *Tree) insertFixup3(n *Node) {
	uncle := n.uncle()
	if t.isRed(uncle) {
		n.parent.color, uncle.color = colorBlock, colorBlock
		n.grandparent().color = colorRed
		t.insertFixup1(n.grandparent())
	} else {
		t.insertFixup4(n)
	}
}

func (t *Tree) insertFixup4(n *Node) {
	grandparent := n.grandparent()
	if n == n.parent.right && n.parent == grandparent.left {
		t.leftRotate(n.parent)
		n = n.left
	} else if n == n.parent.left && n.parent == grandparent.right {
		t.rightRotate(n.parent)
		n = n.right
	}
	t.insertFixup5(n)
}

func (t *Tree) insertFixup5(n *Node) {
	n.parent.color = colorBlock
	grandparent := n.grandparent()
	grandparent.color = colorRed
	if n == n.parent.left && n.parent == grandparent.left {
		t.rightRotate(grandparent)
	} else if n == n.parent.right && n.parent == grandparent.right {
		t.leftRotate(grandparent)
	}
}

func (t *Tree) eraseFixup1(n *Node) {
	if n.parent == nil {
		return
	}
	t.eraseFixup2(n)
}

func (t *Tree) eraseFixup2(n *Node) {
	sibling := n.sibling()
	if t.isRed(sibling) {
		n.parent.color = colorRed
		sibling.color = colorBlock
		if n == n.parent.left {
			t.leftRotate(n.parent)
		} else {
			t.rightRotate(n.parent)
		}
	}
	t.eraseFixup3(n)
}

func (t *Tree) eraseFixup3(n *Node) {
	sibling := n.sibling()
	if !t.isRed(n.parent) &&
		!t.isRed(sibling) &&
		!t.isRed(sibling.left) &&
		!t.isRed(sibling.right) {
		sibling.color = colorRed
		t.eraseFixup1(n.parent)
	} else {
		t.eraseFixup4(n)
	}
}

func (t *Tree) eraseFixup4(n *Node) {
	sibling := n.sibling()
	if t.isRed(n.parent) &&
		!t.isRed(sibling) &&
		!t.isRed(sibling.left) &&
		!t.isRed(sibling.right) {
		sibling.color = colorRed
		n.parent.color = colorBlock
	} else {
		t.eraseFixup5(n)
	}
}

func (t *Tree) eraseFixup5(n *Node) {
	sibling := n.sibling()
	if n == n.parent.left &&
		!t.isRed(sibling) &&
		t.isRed(sibling.left) &&
		!t.isRed(sibling.right) {
		sibling.color = colorRed
		sibling.left.color = colorBlock
		t.rightRotate(sibling)
	} else if n == n.parent.right &&
		!t.isRed(sibling) &&
		t.isRed(sibling.right) &&
		!t.isRed(sibling.left) {
		sibling.color = colorRed
		sibling.right.color = colorBlock
		t.leftRotate(sibling)
	}
	t.eraseFixup6(n)
}

func (t *Tree) eraseFixup6(n *Node) {
	sibling := n.sibling()
	sibling.color = t.isRed(n.parent)
	n.parent.color = colorBlock
	if n == n.parent.left && t.isRed(sibling.right) {
		sibling.right.color = colorBlock
		t.leftRotate(n.parent)
	} else if t.isRed(sibling.left) {
		sibling.left.color = colorBlock
		t.rightRotate(n.parent)
	}
}
