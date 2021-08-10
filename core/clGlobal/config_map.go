package clGlobal

import (
	"github.com/xiaolan580230/clhttp-framework/core/skylog"
	"sync"
)

var mConfigMap map[string] string
var mLocker sync.RWMutex
func init() {
	mConfigMap = make(map[string] string)
}

// 加载配置
func LoadConfig(_section, _key, _default string) string {
	mLocker.Lock()
	defer mLocker.Unlock()

	var temp string
	if conf == nil {
		skylog.LogErr("无法找到配置指针,请先执行Init指定配置文件")
		return ""
	}
	conf.GetStr(_section, _key, _default, &temp)
	mConfigMap[_section + "_" + _key] = temp
	return temp
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