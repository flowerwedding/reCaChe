# Readme

bigCache的分片思想，对LRU的简单优化：

创建有N个分片的数组，每个分片包含自己的带有锁的缓存实例。当需要缓存某个key的记录时，首先通过散列函数`hash(key)%N`选择一个分片，这个过程不需要锁。然后获取具体分片的缓存锁并对缓存进行读写。记录的读取与之类似。N个goroutine每次请求正好平均落在各自的分片上，这样竞争就会减小，即使有多个goroutine落在同一分片上，如果Hash比较平均，则单个分片的压力也会减小。竞争减小后，等待获取锁的时间变短，延迟提高。

平均落在不同的切片/数据不倾斜，可以使用虚拟节点的思想。

其他：[分布式缓存](https://geektutu.com/post/geecache.html)