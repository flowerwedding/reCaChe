/**
 * @Title  consistent
 * @description  一致性哈希
 * @Author  沈来
 * @Update  2020/8/19 18:46
 **/
package consistent

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//依赖注入方法，允许用于替换成自定义的函数
type Hash func(data []byte) uint32

type Map struct {
	hash        Hash
	replicas    int//虚拟节点倍数
	keys        []int//哈希环
	hashMap     map[int]string//key是虚拟节点的值，value是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//添加真实节点
func (m *Map) Add(keys ...string) {
	//对每一个真实节点 key，对应创建 m.replicas 个虚拟节点，通过添加编号的方式区分不同虚拟节点
	for _, key := range keys {
		for i:= 0; i < m.replicas; i++ {
			//使用 m.hash() 计算虚拟节点的哈希值，使用 append(m.keys, hash) 添加到环上
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			//在 hashMap 中增加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	//环上的哈希值排序
	sort.Ints(m.keys)
}

//选择节点
func (m *Map) Get(key string) string{
	if len(m.keys) == 0 {
		return ""
	}

	//计算 key 的哈希值
	hash := int(m.hash([]byte(key)))
	//顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值
	//如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//通过 hashMap 映射得到真实的节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}