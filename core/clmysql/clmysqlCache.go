package clmysql

import (
	"time"
	"sync"
)

// 数据库缓存管理器

type DbCache struct {
	expire time.Time
	data *DbResult
}

type DbCacheMgr struct {
	DataPool map[string] DbCache
	DbLock sync.RWMutex
}


// 创建缓存管理器
func NewCacheMgr() (*DbCacheMgr) {
	var cacheMgr = DbCacheMgr{
		DataPool: make(map[string] DbCache),
	}
	return &cacheMgr
}

// 获取缓存
func (cache *DbCacheMgr) GetCache(key string) (*DbResult) {
	cache.DbLock.RLock()
	defer cache.DbLock.RUnlock()

	if result, ok := cache.DataPool[key]; ok {
		if result.expire.Unix() > time.Now().Unix() {
			return result.data
		}
	}
	return nil
}

// 设置缓存
func (cache *DbCacheMgr)SetCache(key string, data *DbResult, sec uint32) {
	cache.DbLock.Lock()
	defer cache.DbLock.Unlock()

	cache.DataPool[key] = DbCache{
		expire: time.Unix(time.Now().Unix() + int64(sec), 0),
		data: data,
	}
}