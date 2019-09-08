package rbtree

import (
	"errors"
)

const (
	colorRed   = true
	colorBlack = false
)

var (
	//ErrNoData :
	ErrNoData = errors.New("tree has no data")
)

type node struct {
	k, v                interface{}
	parent, left, right *node
	color               bool
}

func newNode(k, v interface{}) *node {
	return &node{k: k, v: v}
}

func rotateLeft(n *node) *node {
	rc := n.right
	n.right = rc.left
	rc.left = n
	rc.color = colorRed
	return rc
}

func rotateRight(n *node) *node {
	lc := n.left
	n.left = lc.right
	lc.right = n
	lc.color = colorRed
	return lc
}

func colorFlip(n *node) *node {
	n.color = !n.color
	if n.left != nil {
		n.left.color = !n.left.color
	}
	if n.right != nil {
		n.right.color = !n.right.color
	}
	return n
}

func fixUp(n *node) *node {
	if n.right.color == colorRed {
		n = rotateLeft(n)
	}
	if n.left.color == colorRed &&
		n.left.left.color == colorRed {
		n = rotateRight(n)
	}
	if n.left.color == colorRed &&
		n.right.color == colorRed {
		n = colorFlip(n)
	}
	return n
}

type tree struct {
	root    *node
	size    int
	compare func(a, b interface{}) int
}

func (t *tree) init(compare func(a, b interface{}) int) {
	t.root = nil
	t.size = 0
	t.compare = compare
}

func (t *tree) insert(n *node, k, v interface{}) (*node, bool) {
	ok := false
	if n == nil {
		return &node{k: k, v: v, color: colorRed}, true
	}
	if t.compare(k, n.k) > 0 {
		n.left, ok = t.insert(n.left, k, v)
	}
}
