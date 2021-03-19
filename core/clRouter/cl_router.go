package clRouter

import (
	"hongxia_api/core/clUserPool"
	"sync"
)

type RouterParam struct {
	Key string				// 参数的key
	Def string				// 参数的默认值
	PType int				// 参数的校验类型
	Static bool				// 是否严格模式
}


type RouterRule struct {
	Ac string					// 路由名称
	Param []RouterParam			// 路由参数
	Static bool					// 是否必须
	Callback func(_uInfo *clUserPool.ClNetUserInfo, _params map[string]string) string			// 回调函数
	Login bool					// 是否需要登录
}

var mRouterMap map[string] RouterRule
var mRouterLock sync.RWMutex

// 初始化
func init() {
	mRouterMap = make(map[string] RouterRule)
}

// 添加路由规则
func AddRule(_info RouterRule) {
	mRouterLock.Lock()
	defer mRouterLock.Unlock()

	mRouterMap[_info.Ac] = _info
}


// 获取路由规则
func GetRule(_ac string) *RouterRule {
	mRouterLock.RLock()
	defer mRouterLock.RUnlock()

	rule, exist := mRouterMap[_ac]
	if !exist {
		return nil
	}
	return &rule
}