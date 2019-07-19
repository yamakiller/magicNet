package aoi

/*const (
	// NodeTypeNormal 非叶节点
	NodeTypeNormal = 0
	// NodeTypeLeaf 叶节点
	NodeTypeLeaf = 1
)

// ChildrenNum 有几个儿子
const ChildrenNum int = 4

// NodeType 节点类型
type NodeType int

// QuadTreeNode : 四叉树节点
type QuadTreeNode struct {
	level     int
	bounds    Rect
	parent    *QuadTreeNode
	nodeType  NodeType
	childrens [ChildrenNum]*QuadTreeNode
	itemCount int
	items     []interface{}

	NodeCapacity int
}

func (qtn *QuadTreeNode) Insert(item interface{}) bool {
	if qtn.nodeType == NodeTypeNormal {
		//LABLE_NORMAL:
	normal_lable:
		if v, ok := item.(Point); ok {
			index := qtn.bounds.GetQuadrantPoint(v) - 1
			if index >= 0 {
				return qtn.childrens[index].Insert(item)
			}
			return false
		}
		panic("item to Point fail")
	}

	if qtn.itemCount < qtn.NodeCapacity {
		if rect, ok := item.(Rect); ok {
			if qtn.bounds.Contains(rect) {
				qtn.itemCount++
				qtn.items = append(qtn.items, item)
				item
			}
		}
	}
}*/
