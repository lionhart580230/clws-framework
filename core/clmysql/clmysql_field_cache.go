package clmysql

import (
	"sync"
	"fmt"
	"time"
)

// 字段缓存
type FieldCache struct {
	expire uint32			// 上次更新时间戳
	fileds []string			// 字段名称列表
}


var FieldsCacheList map[string] FieldCache
var FieldsCacheLocker sync.RWMutex


func init() {
	FieldsCacheList = make(map[string] FieldCache)
}

func (this *DBPointer) GetFields(dbname string, tablename string) []string {
	var cache_key = fmt.Sprintf("%v_%v", dbname, tablename)
	FieldsCacheLocker.RLock()
	cache_val, exists := FieldsCacheList[cache_key]
	FieldsCacheLocker.RUnlock()

	if !exists || cache_val.expire < uint32(time.Now().Unix()) - 300 {
		this.RestoreFieldsCache(tablename)
	}
	return nil
}

func (this *DBPointer) RestoreFieldsCache(tablename string) error {
	resp, err := this.Query("DESC %v", tablename)
	if err != nil {
		return err
	}

	if resp == nil || resp.Length == 0 {
		return nil
	}

	return nil
}