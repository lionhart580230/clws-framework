package clredis

import "sync"

// 内存缓存，为Redis抵挡短时间内一样的请求

type redisCache struct {
	cachelock sync.RWMutex
	cachePool map[string] []byte
}

