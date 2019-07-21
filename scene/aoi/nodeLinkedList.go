package aoi

import (
	"container/list"
	"math"
	"reflect"
	"unsafe"

	"github.com/yamakiller/magicNet/engine/util"
)

const (
	xTLink = 0
	yTLink = 1
)

// 内存控制上需要优化
type nodeLinkedList struct {
	list.List
	skipCount int
	linkType  int
}

func (nll *nodeLinkedList) Insert(obj *Object) {
	if nll.linkType == xTLink {
		nll.insertX(obj)
	} else {
		nll.insertY(obj)
	}
}

func (nll *nodeLinkedList) insertX(obj *Object) {
	if nll.Front() == nil {
		obj.xnode = nll.PushFront(obj)
	} else {
		slowCursor := nll.Front()
		util.AssertEmpty(slowCursor, "insert x front is nil")
		skip := int(math.Ceil(float64(nll.Len()) / float64(nll.skipCount)))
		last := (*Object)(unsafe.Pointer(reflect.ValueOf(nll.Back().Value).Pointer()))
		util.AssertEmpty(last, "not Object")
		if last.position.X > obj.position.X {
			for i := 0; i < nll.skipCount; i++ {
				fastCursor := nll.fastCursor(skip, (*Object)(unsafe.Pointer(reflect.ValueOf(slowCursor.Value).Pointer())))
				fastCursorValue := (*Object)(unsafe.Pointer(reflect.ValueOf(fastCursor.Value).Pointer()))
				util.AssertEmpty(fastCursorValue, "not Object")
				if fastCursorValue.position.X < obj.position.X {
					slowCursor = fastCursor
					continue
				}

				// 慢指针移动到快指针位置
				for slowCursor != nil {
					slowCursorValue := (*Object)(unsafe.Pointer(reflect.ValueOf(slowCursor.Value).Pointer()))
					if slowCursorValue.position.X >= obj.position.X {
						nll.InsertAfter(obj, slowCursor)
						return
					}
					slowCursor = slowCursor.Next()
				}
			}
		}

		if obj.xnode == nil {
			obj.xnode = nll.PushBack(obj)
		}
	}
}

func (nll *nodeLinkedList) insertY(obj *Object) {
	if nll.Front() == nil {
		obj.ynode = nll.PushFront(obj)
	} else {
		slowCursor := nll.Front()
		util.AssertEmpty(slowCursor, "insert y front is nil")
		skip := int(math.Ceil(float64(nll.Len()) / float64(nll.skipCount)))
		last := (*Object)(unsafe.Pointer(reflect.ValueOf(nll.Back().Value).Pointer()))
		util.AssertEmpty(last, "not Object")
		if last.position.Y > obj.position.Y {
			for i := 0; i < nll.skipCount; i++ {
				fastCursor := nll.fastCursor(skip, (*Object)(unsafe.Pointer(reflect.ValueOf(slowCursor.Value).Pointer())))
				fastCursorValue := (*Object)(unsafe.Pointer(reflect.ValueOf(fastCursor.Value).Pointer()))
				util.AssertEmpty(fastCursorValue, "not Object")
				if fastCursorValue.position.Y < obj.position.Y {
					slowCursor = fastCursor
					continue
				}

				// 慢指针移动到快指针位置
				for slowCursor != nil {
					slowCursorValue := (*Object)(unsafe.Pointer(reflect.ValueOf(slowCursor.Value).Pointer()))
					if slowCursorValue.position.Y >= obj.position.Y {
						nll.InsertAfter(obj, slowCursor)
						return
					}
					slowCursor = slowCursor.Next()
				}
			}
		}

		if obj.ynode == nil {
			obj.ynode = nll.PushBack(obj)
		}
	}
}

func (nll *nodeLinkedList) fastCursor(skip int, currObj *Object) *list.Element {
	skipLink := currObj
	switch nll.linkType {
	case xTLink:
		for i := 1; i <= skip; i++ {
			if skipLink.xnode.Next() == nil {
				break
			}
			skipLink = (*Object)(unsafe.Pointer(reflect.ValueOf(skipLink.xnode.Next().Value).Pointer()))
		}
		return skipLink.xnode
	case yTLink:
		for i := 1; i <= skip; i++ {
			if skipLink.ynode.Next() == nil {
				break
			}
			skipLink = (*Object)(unsafe.Pointer(reflect.ValueOf(skipLink.ynode.Next().Value).Pointer()))
		}
		return skipLink.ynode
	default:
		return nil
	}
}
