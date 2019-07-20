package aoi

import (
	"container/list"
	"math"
	"reflect"
	"unsafe"
)

// NewAoiSpace : 创建一个AOI空间
func NewAoiSpace() *Space {
	o := &Space{maps: make(map[uint32]*Object, 64),
		xlink: &nodeLinkedList{skipCount: 10, linkType: xTLink},
		ylink: &nodeLinkedList{skipCount: 10, linkType: yTLink}}
	o.xlink.Init()
	o.ylink.Init()
	return o
}

const (
	setMapSize = 8
)

// Space : Aoi空间
type Space struct {
	maps  map[uint32]*Object
	xlink *nodeLinkedList
	ylink *nodeLinkedList
}

//Enter : 进入
func (sp *Space) Enter(id uint32, x float64, y float64) *Object {
	obj := sp.maps[id]
	if obj != nil {
		return obj
	}

	obj = &Object{id: id,
		moveSet:     objectSet{num: setMapSize, set: make(map[uint32]interface{}, setMapSize)},
		moveOnlySet: objectSet{num: setMapSize, set: make(map[uint32]interface{}, setMapSize)},
		entersSet:   objectSet{num: setMapSize, set: make(map[uint32]interface{}, setMapSize)},
		leavesSet:   objectSet{num: setMapSize, set: make(map[uint32]interface{}, setMapSize)}}

	obj.position.X = x
	obj.position.Y = y

	sp.xlink.Insert(obj)
	sp.ylink.Insert(obj)

	sp.maps[id] = obj
	return obj
}

//Update : 更新
func (sp *Space) Update(id uint32, area Vector2, x float64, y float64) *Object {
	obj := sp.maps[id]
	if obj == nil {
		return nil
	}

	return sp.UpdateObject(obj, area, x, y)
}

// UpdateArea : 更新自己的区域
func (sp *Space) UpdateArea(obj *Object, area Vector2) *Object {
	return sp.UpdateObject(obj, area, obj.position.X, obj.position.Y)
}

// UpdateObject : 目标对象更新
func (sp *Space) UpdateObject(obj *Object, area Vector2, x float64, y float64) *Object {
	obj.moveOnlySet.Rest()
	obj.moveSet.ToCopy(&obj.moveOnlySet)

	sp.Move(obj, x, y)

	sp.findArea(obj, area)

	obj.entersSet.Rest()
	obj.moveSet.Except(&obj.moveOnlySet, &obj.entersSet)

	for k := range obj.entersSet.set {
		if v := sp.maps[k]; v != nil {
			v.moveSet.Add(obj.id)
		}
	}

	obj.leavesSet.Rest()
	obj.moveOnlySet.Except(&obj.moveSet, &obj.leavesSet)

	tmpSet := objectSet{num: 16, set: make(map[uint32]interface{}, 16)}
	obj.moveOnlySet.Except(&obj.entersSet, &tmpSet)
	obj.moveOnlySet.Rest()
	tmpSet.Except(&obj.leavesSet, &obj.moveOnlySet)
	return obj
}

// Move : 移动
func (sp *Space) Move(obj *Object, x float64, y float64) {
	// X
	if math.Abs(float64(obj.position.X-x)) > 0 {
		if x > obj.position.X {
			cur := obj.xnode.Next()
			for cur != nil {

				curValue := (*Object)(unsafe.Pointer(reflect.ValueOf(cur.Value).Pointer()))
				if x < curValue.position.X {
					sp.xlink.Remove(obj.xnode)

					obj.position.X = x

					obj.xnode = sp.xlink.InsertBefore(obj, cur)

					break
				} else if cur.Next() == nil {
					sp.xlink.Remove(obj.xnode)

					obj.position.X = x

					obj.xnode = sp.xlink.InsertAfter(obj, cur)
					break
				}
				cur = cur.Next()
			}
		} else {
			cur := obj.xnode.Prev()

			for cur != nil {
				curValue := (*Object)(unsafe.Pointer(reflect.ValueOf(cur.Value).Pointer()))
				if x > curValue.position.X {
					sp.xlink.Remove(obj.xnode)

					obj.position.X = x

					obj.xnode = sp.xlink.InsertAfter(obj, cur)

					break
				} else if cur.Prev() == nil {

					sp.xlink.Remove(obj.xnode)

					obj.position.X = x

					obj.xnode = sp.xlink.InsertAfter(obj, cur)

					break
				}

				cur = cur.Prev()
			}
		}
	}
	// Y轴
	if math.Abs(float64(obj.position.Y-y)) > 0 {
		if y > obj.position.Y {
			cur := obj.xnode.Next()
			for cur != nil {

				curValue := (*Object)(unsafe.Pointer(reflect.ValueOf(cur.Value).Pointer()))
				if y < curValue.position.Y {
					sp.ylink.Remove(obj.ynode)

					obj.position.Y = y

					obj.ynode = sp.ylink.InsertBefore(obj, cur)

					break
				} else if cur.Next() == nil {
					sp.ylink.Remove(obj.ynode)

					obj.position.Y = y

					obj.ynode = sp.ylink.InsertAfter(obj, cur)
					break
				}
				cur = cur.Next()
			}
		} else {
			cur := obj.ynode.Prev()

			for cur != nil {
				curValue := (*Object)(unsafe.Pointer(reflect.ValueOf(cur.Value).Pointer()))
				if y > curValue.position.Y {
					sp.ylink.Remove(obj.ynode)

					obj.position.Y = y

					obj.ynode = sp.ylink.InsertAfter(obj, cur)

					break
				} else if cur.Prev() == nil {

					sp.ylink.Remove(obj.ynode)

					obj.position.Y = y

					obj.ynode = sp.ylink.InsertAfter(obj, cur)

					break
				}

				cur = cur.Prev()
			}
		}
	}

	obj.position.X = x
	obj.position.Y = y
}

// LeaveNode : 离开
func (sp *Space) LeaveNode(id uint32) []uint32 {
	obj := sp.maps[id]
	if obj == nil {
		return nil
	}

	sp.xlink.Remove(obj.xnode)

	sp.ylink.Remove(obj.ynode)

	delete(sp.maps, id)

	nodes := obj.moveSet.ToArray()

	obj.xnode = nil
	obj.ynode = nil

	return nodes
}

func (sp *Space) find(id uint32, area Vector2) *Object {

	obj := sp.maps[id]
	if obj == nil {
		return nil
	}

	sp.findArea(obj, area)
	return obj
}

func (sp *Space) findArea(obj *Object, area Vector2) {
	obj.moveSet.Rest()

	for i := 0; i < 2; i++ {
		var cur *list.Element
		if i == 0 {
			cur = obj.xnode.Next()
		} else {
			cur = obj.xnode.Prev()
		}

		for cur != nil {
			curValue := (*Object)(unsafe.Pointer(reflect.ValueOf(cur.Value).Pointer()))
			if math.Abs(math.Abs(curValue.position.X)-math.Abs(obj.position.X)) > area.X {
				break
			} else if math.Abs(math.Abs(curValue.position.Y)-math.Abs(obj.position.Y)) > area.Y {
				if Distance2(obj.position, curValue.position) <= area.X {
					obj.moveSet.Add(curValue.id)
				}
			}

			if i == 0 {
				cur = obj.xnode.Next()
			} else {
				cur = obj.xnode.Prev()
			}
		}
	}

	for i := 0; i < 2; i++ {
		var cur *list.Element
		if i == 0 {
			cur = obj.ynode.Next()
		} else {
			cur = obj.ynode.Prev()
		}

		for cur != nil {
			curValue := (*Object)(unsafe.Pointer(reflect.ValueOf(cur.Value).Pointer()))
			if math.Abs(math.Abs(curValue.position.Y)-math.Abs(obj.position.Y)) > area.Y {
				break
			} else if math.Abs(math.Abs(curValue.position.X)-math.Abs(obj.position.X)) > area.X {
				if Distance2(obj.position, curValue.position) <= area.Y {
					obj.moveSet.Add(curValue.id)
				}
			}

			if i == 0 {
				cur = obj.ynode.Next()
			} else {
				cur = obj.ynode.Prev()
			}
		}
	}
}
