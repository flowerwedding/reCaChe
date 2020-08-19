/**
 * @Title  hasher
 * @description  该算法采用位运算方式在栈上运算，避免堆上内存分配
 * @Author  沈来
 * @Update  2020/8/18 16:32
 **/
package fast

func newDefaultHasher() fnv64a {
	return fnv64a{}
}

type fnv64a struct{}

//fnv哈希算法64位时固定的参数
const (
	offset64 = 14695981039346656037
	prime64 = 1099511628211
)

func (f fnv64a) Sum64 (key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key);i ++{
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash
}