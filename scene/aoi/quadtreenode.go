package aoi

import (
	"magicNet/engine/util"
	"unsafe"
)

const (
	// NodeTypeNormal 非叶节点
	NodeTypeNormal = 0
	// NodeTypeLeaf 叶节点
	NodeTypeLeaf = 1
)

// ChildrenNum 有几个儿子
const ChildrenNum int = 4

// NodeType 节点类型
type NodeType int

type QuadItem struct {
	Node     unsafe.Pointer
	ItemNext unsafe.Pointer
}

type QuadTreeNode struct {
	level     int
	bounds    Rect
	parent    *QuadTreeNode
	nodeType  NodeType
	childrens [ChildrenNum]*QuadTreeNode
	itemCount int
	items     unsafe.Pointer

	NodeCapacity int
	LevelLimit   int
}

func (qtn *QuadTreeNode) Insert(item unsafe.Pointer) bool {
	isCheck := true
	for {
		if qtn.nodeType == NodeTypeNormal || !isCheck {
			isCheck = true
			point := (*Point)(item)
			util.AssertEmpty(point, "QuadTreeNode Insert item not Point")
			index := qtn.bounds.GetQuadrant(point) - 1
			if index >= 0 {
				return qtn.childrens[index].Insert(item)
			}

			return false
		}

		if qtn.itemCount > qtn.NodeCapacity {
			if qtn.bounds.Contains(item) {
				qtn.itemCount++
				//qtn.items.
				pitem := (*QuadItem)(item)
				util.AssertEmpty(pitem, "QuadTreeNode Insert item not QuadItem")
				pitem.ItemNext = qtn.items
				qtn.items = item
				pitem.Node = unsafe.Pointer(qtn)
				return true
			}

			return false
		}

		if qtn.level+1 >= qtn.LevelLimit {
			return false
		}

		qtn.split()
		isCheck = false
	}
}

func (qtn *QuadTreeNode) split() {
	util.Assert(qtn.nodeType == NodeTypeLeaf, "QuadTreeNode split not NodeTypeLeaf")
	qtn.nodeType = NodeTypeNormal

	rect0 := GenRect(qtn.bounds.MidX(), qtn.bounds.Right(), qtn.bounds.MidY(), qtn.bounds.Top())

	rect1 := GenRect(qtn.bounds.Left(), qtn.bounds.MidX(), qtn.bounds.MidY(), qtn.bounds.Top())

	// 第三象限，左下
	rect2 := GenRect(qtn.bounds.Left(), qtn.bounds.MidX(), qtn.bounds.Bottom(), qtn.bounds.MidY())

	// 第四象限，右下
	rect3 := GenRect(qtn.bounds.MidX(), qtn.bounds.Right(), qtn.bounds.Bottom(), qtn.bounds.MidY())

	//qtn.childrens[0] = mAlloc->New(mLevel + 1, mAlloc, NodeTypeLeaf, this, rect0)
	//qtn.childrens[1] = mAlloc->New(mLevel + 1, mAlloc, NodeTypeLeaf, this, rect1)
	//qtn.childrens[2] = mAlloc->New(mLevel + 1, mAlloc, NodeTypeLeaf, this, rect2)
	//qtn.childrens[3] = mAlloc->New(mLevel + 1, mAlloc, NodeTypeLeaf, this, rect3)

	/*for (TItem* it = mItems; it;)
	  {
	      TItem* head = (TItem*)(it->mItemNext);
	      int index = mBounds.GetQuadrant2(it) - 1;
	      mChildrens[index]->Insert(it);
	      it = head;
	  }
	  mItemCount = 0;
	  mItems = nullptr;*/

}
