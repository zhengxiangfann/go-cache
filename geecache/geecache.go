package geecache

import (
	"fmt"
	"sync"
)

// 定义接口 Getter 和 回调函数 Get(key string)([]byte, error)，参数是 key，返回值是 []byte
// 定义函数类型 GetterFunc，并实现 Getter 接口的 Get 方法。

// 函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，
// 也能够传入实现了该接口的结构体作为参数。

type (
	Getter interface {
		Get(key string) ([]byte, error)
	}

	GetterFunc func(key string) ([]byte, error)
)

// 定义几个方法类型, 这个类型实现的接口的方法，在这个函数内部调用 函数类型自己
// 把普通的 函数转换为接口
/*
定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。
这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。

这个语法并没有什么特别的。你其实可以反过来想一想，如果不提供这个把函数转换为接口的函数，
你调用时就需要创建一个struct，然后实现对应的接口，
创建一个实例作为参数，相比这种方式就麻烦得多了。
*/
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type (
	Group struct {
		name      string
		getter    Getter
		mainCache cache
	}
)

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	g := groups[name]
	mu.Unlock()

	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 从 mainCache 中查找缓存，如果存在则返回缓存值。
	if v, exist := g.mainCache.get(key); exist {
		return v, nil
	}

	// 缓存不存在，则调用 load 方法，load 调用 getLocally
	//（分布式场景下会调用 getFromPeer 从其他节点获取），
	//getLocally 调用用户回调函数
	//g.getter.Get() 获取源数据，
	//并且将源数据添加到缓存 mainCache 中（通过 populateCache 方法）
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}

	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
