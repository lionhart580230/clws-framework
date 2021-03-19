package clUserPool

import (
	"github.com/gorilla/websocket"
	"hongxia_api/core/clCommon"
	"strings"
	"sync"
)

// 管理连线用户池
// 当用户成功创建连线的时候，用户池中会对相应的连线创建一个用户对象
// 默认此用户对象是不存在任何状态和数据的，外部可通过一些接口对用户对象进行操作
// 此用户对象默认是常驻内存，可通过修改配置实现数据实时同步到redis中，以便重启后恢复数据


type ClNetUserInfo struct {
	ConnId uint64						// 用户的连线ID
	Conn *websocket.Conn				// 用户的连线
	IsLogin bool						// 是否登录
	IP string							// IP
	Port uint32							// 端口
	Params map[string]string			// 一些可自由设置的扩展参数
	ParamLock sync.RWMutex				// 参数锁
}


var mUserPoolMap map[uint64] *ClNetUserInfo
var mGlobalId uint64 = 1
var mUserPoolLocker sync.RWMutex

func init() {
	mGlobalId = 1
	mUserPoolMap = make(map[uint64] *ClNetUserInfo)
}

func AddNewUser(_conn *websocket.Conn, _islogin bool) *ClNetUserInfo {
	mUserPoolLocker.Lock()
	defer mUserPoolLocker.Unlock()

	mGlobalId++
	var addr = strings.Split(_conn.LocalAddr().String(), ":")

	var uInfo = &ClNetUserInfo{
		ConnId:    mGlobalId,
		Conn:      _conn,
		IsLogin:   _islogin,
		IP:        addr[0],
		Port:      clCommon.Uint32(addr[1], 0),
		Params:    make(map[string] string),
	}

	mUserPoolMap[mGlobalId] = uInfo

	return uInfo
}