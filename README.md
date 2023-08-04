###### 1.  lru缓存淘汰策略
###### 2.  单机并发缓存
    
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

###### 3.  http服务端
###### 4.  一致性hash
###### 5.  分布式节点
