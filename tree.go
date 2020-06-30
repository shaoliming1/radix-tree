package radix_tree

// 设计：
// 1. 不需要字段标识一个node是不是叶子节点， 如果`children`为 nil, 那么
//    改节点就是叶子节点。
// 2. key应该是在edge上，但是由于一个节点最多只有一个父节点，那么可以每条边都连着
//    一个子节点，所以可以放在子节点上。

// 修正：
// 存在一个问题，就是前缀重合的情况，"abc"和"abcd"形成的radix tree 如图所示
//                root
//               /
//              abc
//              /
//             d
// 可以看出"abc"节点的children不在为空，如果这时我们查找"abc",这颗树会找不到。
// 解决方案： 增加一个结尾符，还是上面这个例子，插入"abc"和"abcd"的树变为：
//               root
//              /
//              abc
//              / \
//            "$"  "d$"
type node struct {
	key      string
	value    interface{}
	children []*node
}

type RadixTree struct {
	root *node
}

// 构造函数
func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: &node{},
	}
}

func longestCommonPrefix(str1 string, str2 string) int {
	i := 0
	for ; i < len(str1) && i < len(str2) && str1[i] == str2[i]; i++ {
	}
	return i
}

// 获得和`remain`有相同前缀的child,如果不存在，返回nil
func getChildWithCommonPrefix(children []*node, remain string) (int, *node) {
	for i, child := range children {
		if len(child.key) > 0 && child.key[0] == remain[0] {
			return i, child
		}
	}
	return -1, nil
}

// ---------- 对外接口 ---------------
func (t *RadixTree) Find(key string) interface{} {
	parent := t.root
	remain := key + "$" // $是结尾符
	// 如果还有未匹配的字符 且 parent.children不为空
	for parent.children != nil && len(remain) > 0 {
		_,child := getChildWithCommonPrefix(parent.children, remain)
		if child != nil{
			idx := longestCommonPrefix(remain, child.key)
			if child.key == remain[:idx]{
				remain = remain[idx:]
				parent = child
			}else{
				return nil
			}
		}
	}
	// 找到了
	if parent.children == nil && len(remain) == 0 {
		return parent.value
	}

	return nil
}

// `Insert` can update the value if key already exists
func (t *RadixTree) Insert(key string, value interface{}) {
	parent := t.root
	remain := key + "$"
	for len(remain) > 0 {
		index, child := getChildWithCommonPrefix(parent.children, remain)
		if child != nil {
			idx := longestCommonPrefix(remain, child.key)
			if child.key == remain[idx:] {
				remain = remain[idx:]
				parent = child
			} else { //部分匹配
				prefix := remain[:idx]
				//
				//parent.children[i]= parent.children[len(parent.children)]
				//child.key = prefix

				// 新建节点
				node1 := &node{
					key: prefix,
				}
				node2 := &node{
					key:   remain[idx:],
					value: value,
				}
				child.key = child.key[idx:]
				//child.value = nil
				parent.children[index] = node1
				node1.children = append(node1.children, child, node2)
				remain = ""
			}
		}else{
			node := &node{
				key:   remain,
				value: value,
			}
			parent.children = append(parent.children, node)
			remain = ""
		}
	}
}

func (t *RadixTree) Delete(key string) {
	remain := key + "$"
	parent := t.root
	var pp *node
	for len(remain) >0{
		index, child := getChildWithCommonPrefix(parent.children, remain)
		if child!=nil {
			idx := longestCommonPrefix(remain, child.key)
			if child.key == remain[:idx] {
				remain = remain[idx:]
				// match
				if len(remain) == 0 && len(child.children)==0{
					if pp!=nil && len(parent.children)==2 {
						 ppindex,_ := getChildWithCommonPrefix(pp.children, parent.key)
						 // delete child
						 parent.children[index] = parent.children[len(parent.children)-1]
						 parent.children = parent.children[:len(parent.children)-1]

						 // 更新另外一个孩子节点的key
						 parent.children[0].key = parent.key + parent.children[0].key

						 // pp成为孩子节点的父节点
						 pp.children[ppindex] = parent.children[0]
					}
				}
				pp = parent
				parent = child
			} else {
				return
			}
		}else{
			return
		}
	}


}
