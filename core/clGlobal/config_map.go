package clGlobal

import (
	"sync"
)

var mConfigMap map[string] string
var mLocker sync.RWMutex
func init() {
	mConfigMap = make(map[string] string)
}



// 获取配置
func GetConfig(_section, _key string) string {
	mLocker.RLock()
	defer mLocker.RUnlock()

	val, exists := mConfigMap[_section + "_" + _key]
	if !exists {
		return ""
	}
	return val
}


// 强制指定配置
func SetConfig(_section, _key, _val string) {
	mLocker.Lock()
	defer mLocker.Unlock()

	mConfigMap[_section + "_" + _key] = _val
}