package lru

import (
	"container/list"
)

type (
	Cache struct {
		maxBytes  int64      //允许使用的最大内存
		nBytes    int64      //当前的内存
		ll        *list.List //
		cache     map[string]*list.Element
		OnEvicted func(key string, value Value) // 是某条记录被移除时的回调函数，可以为 nil。
	}

	// 为了通用性，我们允许值是实现了 Value 接口的任意类型，
	// 该接口只包含了一个方法 Len() int，
	// 用于返回值所占用的内存大小。

	Value interface {
		Len() int
	}

	// 键值对 entry 是双向链表节点的数据类型，
	// 在链表中仍保存每个值对应的 key 的好处在于，
	// 淘汰队首节点时，需要用 key 从字典中删除对应的映射
	entry struct {
		key   string
		value Value
	}
)

// New is the constructor cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
// 如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值。
// c.ll.MoveToFront(ele)，即将链表中的节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的，
// 在这里约定 front 为队尾）
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) // 将节点移动到队头
		kv := ele.Value.(*entry)
		//value, ok = kv.value, true
		return kv.value, true
	}
	return
}

// 这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // c.ll.Back() 取到队首节点，从链表中删除。
	if ele != nil {
		c.ll.Remove(ele) //
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 从字典中 c.cache 删除该节点的映射关系。
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新当前所用的内存 c.nBytes。
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) // 回调删除的动作
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 则更新对应节点的值，并将该节点移到队尾。
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {

		ele1 := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele1
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}
