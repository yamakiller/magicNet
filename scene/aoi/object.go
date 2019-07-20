package aoi

import "container/list"

// Object aoi 对象
type Object struct {
	id       uint32
	position Vector2

	moveSet     objectSet
	moveOnlySet objectSet
	entersSet   objectSet
	leavesSet   objectSet

	xnode *list.Element
	ynode *list.Element
}
