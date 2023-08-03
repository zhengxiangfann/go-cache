package geecache

type (
	ByteView struct {
		b []byte
	}
)

// 实现了Value的接口(ByteView 就是一个)
func (b ByteView) Len() int {
	return len(b.b)
}

// ByteSlice
func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}

// cloneBytes 克隆一个字符数组
func cloneBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
