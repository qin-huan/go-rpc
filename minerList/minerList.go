package minerList

type MinerNode struct {
	Data interface{}
	Next *MinerNode
}

func New (data interface{}) *MinerNode {
	root := MinerNode{
		Data: data,
		Next: nil,
	}
	root.Next = &root
	return &root
}

func (this *MinerNode) IsEmpty (root *MinerNode) bool {
	return root == nil
}

func (this *MinerNode) Find (data interface{}, root *MinerNode) *MinerNode {
	if root == nil {
		return nil
	}
	node := root
	for {
		if node.Data == data {
			break
		}
		node = node.Next
		if node == root {
			return nil
		}
	}
	return node
}

func (this *MinerNode) Length (root *MinerNode) int {
	if root == nil {
		return 0
	}
	head := root.Next
	index := 1
	for {
		if head == root {
			break
		}
		head = head.Next
		index++
	}
	return index
}

func (this *MinerNode) Insert (data interface{}, position *MinerNode) {
	tmp := MinerNode{
		Data: data,
		Next: position.Next,
	}
	position.Next = &tmp
}

func (this *MinerNode) Delete (data interface{}, root *MinerNode) {
	node := this.Find(data, root)

	if node == nil || node.Length(node) == 1 {
		node = nil
	} else {
		node.Data = node.Next.Data
		node.Next = node.Next.Next
	}
}