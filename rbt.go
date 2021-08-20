package rbt

type RBT interface {
	Set(key int, value interface{})
	Del(key int)
	Get(key int) (value interface{}, ok bool)
	Print() (keyList []int, valueList []interface{})
}

const (
	RED   = true
	BLACK = false
)

type node struct {
	parent, right, left *node
	color               bool
	key                 int
	value               interface{}
}

type rbt struct {
	root *node
	null *node
	len  int
}

func (r *rbt) leftRotate(x *node) {

	y := x.right
	x.right = y.left
	if y.left != r.null {
		y.left.parent = x
	}
	y.parent = x.parent
	switch {
	case x.parent == r.null:
		r.root = y
	case x == x.parent.left:
		x.parent.left = y
	default:
		x.parent.right = y
	}

	y.left = x
	x.parent = y

}

func (r *rbt) rightRotate(y *node) {

	x := y.left
	y.left = x.right
	if x.right != r.null {
		x.right.parent = y
	}
	x.parent = y.parent
	switch {
	case y.parent == r.null:
		r.root = x
	case y == y.parent.left:
		y.parent.left = x
	default:
		y.parent.right = x
	}
	x.right = y
	y.parent = x
}

func (r *rbt) insert(z *node) {

	y := r.null
	x := r.root

	for x != r.null {
		y = x
		switch {
		case z.key < x.key:
			x = x.left
		case z.key > x.key:
			x = x.right
		case z.key == x.key:
			x.value = z.value
			return
		}
	}

	z.parent = y
	switch {
	case y == r.null:
		r.root = z
	case z.key < y.key:
		y.left = z
	default:
		y.right = z
	}

	z.left = r.null
	z.right = r.null
	z.color = RED

	r.insertFixUp(z)
	r.len++

}

func (r *rbt) insertFixUp(z *node) {
	for z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			//父节点为左节点
			y := z.parent.parent.right // 叔节点
			if y.color == RED {
				// case 1: 叔节点为红色
				y.color = BLACK
				z.parent.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					//case 2: 叔节点为黑色&当前节点为右节点
					z = z.parent
					r.leftRotate(z)
				}
				//case 3: 叔节点为黑色&当前节点为左节点
				z.parent.color = BLACK
				z.parent.parent.color = RED
				r.rightRotate(z.parent.parent)
			}
		} else {
			//父节点为右节点
			y := z.parent.parent.left // 叔节点
			if y.color == RED {
				// case 1: 叔节点为红色
				y.color = BLACK
				z.parent.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					//case 2: 叔节点为黑色&当前节点为左节点
					z = z.parent
					r.rightRotate(z)
				}
				//case 3: 叔节点为黑色&当前节点为右节点
				z.parent.color = BLACK
				z.parent.parent.color = RED
				r.leftRotate(z.parent.parent)
			}
		}
	}
	r.root.color = BLACK
}

func (r *rbt) transplant(u, v *node) {

	switch {
	case u.parent == r.null:
		//u:根节点
		r.root = v
	case u == u.parent.left:
		//u:左节点
		u.parent.left = v
	default:
		//u:右节点
		u.parent.right = v
	}
	v.parent = u.parent
}

func (r *rbt) minimum(x *node) *node {
	for x.left != r.null {
		x = x.left
	}
	return x
}

func (r *rbt) deleteFixUp(x *node) {

	for x != r.root && x.color == BLACK {
		if x == x.parent.left {
			//x 为左节点
			w := x.parent.right //x 的兄弟节点

			if w.color == RED {
				//case 1:x的兄弟节点w为红色
				w.color = BLACK
				x.parent.color = RED
				r.leftRotate(x.parent)
				w = x.parent.right
			}

			if w.left.color == BLACK && w.right.color == BLACK {
				// case 2:x的兄弟节点w为黑色 & x的两个侄子节点为黑色
				w.color = RED
				x = x.parent
			} else {
				if w.right.color == BLACK {
					// case 3:x的兄弟节点w为黑色 & x的右侄子节点为黑色 & x的左侄子节点为红色
					w.left.color = BLACK
					w.color = RED
					r.rightRotate(w)
					w = x.parent.right
				}

				// case 4:x的兄弟节点w为黑色 & x的右侄子节点为红色
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				r.leftRotate(x.parent)
				x = r.root
			}
		} else {
			//x 为右节点
			w := x.parent.left //x 的兄弟节点

			if w.color == RED {
				//case 1:x的兄弟节点w为红色
				w.color = BLACK
				x.parent.color = RED
				r.rightRotate(x.parent)
				w = x.parent.left
			}

			if w.left.color == BLACK && w.right.color == BLACK {
				// case 2:x的兄弟节点w为黑色 & x的两个侄子节点为黑色
				w.color = RED
				x = x.parent
			} else {
				if w.left.color == BLACK {
					// case 3:x的兄弟节点w为黑色 & x的左侄子节点为黑色 & x的右侄子节点为红色
					w.right.color = BLACK
					w.color = RED
					r.leftRotate(w)
					w = x.parent.left
				}

				// case 4:x的兄弟节点w为黑色 & x的左侄子节点为红色
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK
				r.rightRotate(x.parent)
				x = r.root
			}
		}
	}
	x.color = BLACK
}

func (r *rbt) delete(z *node) {

	y := z
	x := r.null
	yOriginalColor := y.color

	switch {
	case z.left == r.null:
		//z的左节点为空
		x = z.right
		r.transplant(z, z.right)
	case z.right == r.null:
		//z的右节点为空
		x = z.left
		r.transplant(z, z.left)
	default:
		//z的两个子节点都不为空
		y = r.minimum(z.right) //最小节点
		yOriginalColor = y.color
		x = y.right

		if y.parent == z {
			//只对哨兵(x)有作用
			x.parent = y
		} else {
			r.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}

		r.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}

	if yOriginalColor == BLACK {
		//删除的节点为黑色
		// x 的兄弟节点不为空
		r.deleteFixUp(x)
	}

	r.len--
}

func (r *rbt) get(key int) (ret *node) {

	ret = r.null
	x := r.root

	for x != r.null {
		switch {
		case key < x.key:
			x = x.left
		case key > x.key:
			x = x.right
		case key == x.key:
			ret = x
			x = r.null
		}
	}

	return
}

func (r *rbt) set(key int, value interface{}) {
	r.insert(&node{
		key:   key,
		value: value,
	})
}

func (r *rbt) del(key int) {

	x := r.get(key)
	if x == r.null {
		return
	}

	r.delete(x)
}

func (r *rbt) print() (keyList []int, valueList []interface{}) {

	stack := make([]*node, r.len, r.len)
	stackIndex := -1
	x := r.root

	for x != r.null || stackIndex > -1 {

		for x != r.null {
			stackIndex++
			stack[stackIndex] = x
			x = x.left
		}

		if stackIndex > -1 {
			x = stack[stackIndex]
			stackIndex--

			keyList = append(keyList, x.key)
			valueList = append(valueList, x.value)

			x = x.right
		}

	}

	return
}

func GenRBT() RBT {
	null := &node{color: BLACK}
	return &rbt{
		root: null,
		null: null,
	}
}

func (r *rbt) Get(key int) (value interface{}, ok bool) {

	if x := r.get(key); x != r.null {
		value, ok = x.value, true
	}

	return
}

func (r *rbt) Set(key int, value interface{}) {
	r.set(key, value)
}

func (r *rbt) Del(key int) {
	r.del(key)
}

func (r *rbt) Print() (keyList []int, valueList []interface{}) {
	return r.print()
}
