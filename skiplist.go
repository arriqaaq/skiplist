/*
	https://www.epaperpress.com/sortsearch/download/skiplist.pdf
*/

package skiplist

import (
	"bytes"
	"math/rand"
)

const (
	SKIPLIST_MAXLEVEL    = 100     /* For 2^100 elements */
	SKIPLIST_Probability = 1.0 / 4 /* Skiplist probability = 1/4 */
)

type (
	sklLevel struct {
		forward *Node
		span    uint64
	}

	// (hop table)
	Node struct {
		key   string
		value interface{}
		level []*sklLevel
	}

	// Node in skip list (jump table)
	Skiplist struct {
		head   *Node
		tail   *Node
		length int64
		level  int
	}
)

// Returns a random level for the new skiplist node we are going to create.
// The return value of this function is between 1 and SKIPLIST_MAXLEVEL
// (both inclusive), with a powerlaw-alike distribution where higher
// levels are less likely to be returned.
func randomLevel() int {
	level := 1
	for float64(rand.Int31()&0xFFFF) < float64(SKIPLIST_Probability*0xFFFF) {
		level += 1
	}
	if level < SKIPLIST_MAXLEVEL {
		return level
	}

	return SKIPLIST_MAXLEVEL
}

func createNode(level int, key string, value interface{}) *Node {
	node := &Node{
		key:   key,
		value: value,
		level: make([]*sklLevel, level),
	}

	for i := range node.level {
		node.level[i] = new(sklLevel)
	}

	return node
}

func (n *Node) Key() string {
	return n.key
}

func (n *Node) Value() interface{} {
	return n.value
}

func New() *Skiplist {
	return &Skiplist{
		level: 1,
		head:  createNode(SKIPLIST_MAXLEVEL, "", nil),
	}
}

func (z *Skiplist) exists(key string) (*Node, bool) {
	x := z.head
	for i := z.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			bytes.Compare([]byte(x.level[i].forward.key), []byte(key)) <= 0 {
			x = x.level[i].forward
		}

		if x.key == key {
			return x, true
		}
	}

	return nil, false
}

func (z *Skiplist) Update(key string, value interface{}) {
	x := z.head
	for i := z.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			bytes.Compare([]byte(x.level[i].forward.key), []byte(key)) <= 0 {
			x = x.level[i].forward
		}

		if x.key == key {
			x.value = value
		}
	}
}

/*
	Insert a new node in the skiplist.
*/
func (z *Skiplist) Set(key string, value interface{}) *Node {
	/*

		https://www.youtube.com/watch?v=UGaOXaXAM5M
		https://www.youtube.com/watch?v=NDGpsfwAaqo

		The update array stores previous pointers for each level, new node
		will be added after them. rank array stores the rank value of each skiplist node.

		Steps:

		generate update and rank array
		create a new node with random level
		insert new node according to update and rank info
		update other necessary infos, such as span, length.
	*/

	if n, exists := z.exists(key); exists {
		n.value = value
		return n
	}

	updates := make([]*Node, SKIPLIST_MAXLEVEL)
	rank := make([]uint64, SKIPLIST_MAXLEVEL)

	x := z.head
	for i := z.level - 1; i >= 0; i-- {
		/* store rank that is crossed to reach the insert position */
		if i == z.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		if x.level[i] != nil {
			for x.level[i].forward != nil &&
				bytes.Compare([]byte(x.level[i].forward.key), []byte(key)) <= 0 {
				rank[i] += x.level[i].span
				x = x.level[i].forward
			}
		}
		updates[i] = x
	}

	level := randomLevel()
	if level > z.level { // add a new level
		for i := z.level; i < level; i++ {
			rank[i] = 0
			updates[i] = z.head
			updates[i].level[i].span = uint64(z.length)
		}
		z.level = level
	}

	x = createNode(level, key, value)
	for i := 0; i < level; i++ {
		x.level[i].forward = updates[i].level[i].forward
		updates[i].level[i].forward = x

		/* update span covered by update[i] as x is inserted here */
		x.level[i].span = updates[i].level[i].span - (rank[0] - rank[i])
		updates[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	/* increment span for untouched levels */
	for i := level; i < z.level; i++ {
		updates[i].level[i].span++
	}

	if x.level[0].forward == nil {
		z.tail = x
	}

	z.length++
	return x
}

func (z *Skiplist) deleteNode(x *Node, updates []*Node) {
	for i := 0; i < z.level; i++ {
		if updates[i].level[i].forward == x {
			updates[i].level[i].span += x.level[i].span - 1
			updates[i].level[i].forward = x.level[i].forward
		} else {
			updates[i].level[i].span--
		}
	}

	for z.level > 1 && z.head.level[z.level-1].forward == nil {
		z.level--
	}

	z.length--
}

/* Delete an element with matching key from the skiplist. */
func (z *Skiplist) Delete(key string) {
	update := make([]*Node, SKIPLIST_MAXLEVEL)

	x := z.head
	for i := z.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			bytes.Compare([]byte(x.level[i].forward.key), []byte(key)) < 0 {
			x = x.level[i].forward
		}
		update[i] = x
	}

	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and object. */
	x = x.level[0].forward
	if x != nil && x.key == key {
		z.deleteNode(x, update)
		return
	}
}

/* Get an element with matching key from the skiplist. */
func (z *Skiplist) Get(key string) *Node {
	x := z.head
	for i := z.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			bytes.Compare([]byte(x.level[i].forward.key), []byte(key)) <= 0 {
			x = x.level[i].forward
		}

		if x.key == key {
			return x
		}
	}

	return nil
}

func (z *Skiplist) Keys() []string {
	keys := make([]string, 0, 1)

	x := z.head
	for x.level[0].forward != nil {
		keys = append(keys, x.level[0].forward.key)
		x = x.level[0].forward
	}

	return keys
}
