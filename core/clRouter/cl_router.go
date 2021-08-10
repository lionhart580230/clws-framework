package clRouter

import (
	"github.com/gorilla/websocket"
	"github.com/xiaolan580230/clws-framework/core/clDebug"
	"github.com/xiaolan580230/clws-framework/core/clPacket"
	"github.com/xiaolan580230/clws-framework/core/clUserPool"
	"sync"
)

type RouterParam struct {
	Key string				// 参数的key
	Def string				// 参数的默认值
	PType int				// 参数的校验类型
	Static bool				// 是否严格模式
}




func JCode(_rc string, _param string, _data interface{}) *clPacket.RuleCBResp {
	return &clPacket.RuleCBResp{
		RC:  _rc,
		Param: _param,
		Data:  _data,
	}
}

type RouterRule struct {
	Ac string					// 路由名称
	Param []RouterParam			// 路由参数
	Callback func(_uInfo *clUserPool.ClNetUserInfo, _params map[string]string) *clPacket.RuleCBResp			// 回调函数
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


// 发送消息
func SendMessage(_user *clUserPool.ClNetUserInfo, _ack uint32, _rc string, _param string, _data interface{}) {
	_user.ConnLock.Lock()
	defer _user.ConnLock.Unlock()

	var p = clPacket.NewPacketResp(_ack, JCode(_rc, _param, _data))
	err := _user.Conn.WriteMessage(websocket.TextMessage, []byte(p))
	if err != nil {
		clDebug.Err("发送消息失败! 错误:%v", err)
	}
}