package clUserPool

import (
	"github.com/gorilla/websocket"
	"github.com/xiaolan580230/clws-framework/core/clCommon"
	"strings"
	"sync"
	"time"
)

// 管理连线用户池
// 当用户成功创建连线的时候，用户池中会对相应的连线创建一个用户对象
// 默认此用户对象是不存在任何状态和数据的，外部可通过一些接口对用户对象进行操作
// 此用户对象默认是常驻内存，可通过修改配置实现数据实时同步到redis中，以便重启后恢复数据


type ClNetUserInfo struct {
	ConnId uint64						// 用户的连线ID
	Flags uint64						// 用户标示
	Conn *websocket.Conn				// 用户的连线
	IsLogin bool						// 是否登录
	LogoutTime uint32					// 离线时间
	IP string							// IP
	Port uint32							// 端口
	Params map[string]string			// 一些可自由设置的扩展参数
	ParamLock sync.RWMutex				// 参数锁
	ConnLock sync.RWMutex				// 连线消息锁
	Token string						// 登录密钥
}


var mUserPoolMap map[uint64] *ClNetUserInfo
var mGlobalId uint64 = 1
var mUserPoolLocker sync.RWMutex

func init() {
	mGlobalId = 1
	mUserPoolMap = make(map[uint64] *ClNetUserInfo)
}


// 添加新用户
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

// 移除老用户
func RemoveUser(_id uint64) {
	mUserPoolLocker.Lock()
	defer mUserPoolLocker.Unlock()

	obj, exist := mUserPoolMap[_id]
	if !exist {
		return
	}
	obj.IsLogin = false
	obj.LogoutTime = uint32(time.Now().Unix())
}

// 根据用户连线Id获取用户指针
func GetUserById(_id uint64) *ClNetUserInfo {
	mUserPoolLocker.RLock()
	defer mUserPoolLocker.RUnlock()
	return mUserPoolMap[_id]
}


// 根据用户flag获取用户指针
func GetUserByFlags(_flag uint64) *ClNetUserInfo {
	mUserPoolLocker.RLock()
	defer mUserPoolLocker.RUnlock()

	for _, val := range mUserPoolMap {
		if val.Flags == _flag && val.IsLogin {
			return val
		}
	}
	return nil
}


// 根据用户指定属性获取用户指针
func GetUserByParams(_key, _val string) *ClNetUserInfo {
	mUserPoolLocker.RLock()
	defer mUserPoolLocker.RUnlock()

	for _, val := range mUserPoolMap {
		if val.Params[_key] == _val {
			return val
		}
	}
	return nil
}


// 根据用户指定属性获取用户列表
func GetUsersByParams(_key, _val string) []*ClNetUserInfo {
	mUserPoolLocker.RLock()
	defer mUserPoolLocker.RUnlock()
	var userList = make([]*ClNetUserInfo, 0)
	for _, val := range mUserPoolMap {
		if val.Params[_key] == _val {
			userList = append(userList, val)
		}
	}
	return userList
}


// 自动清理5分钟内离线的用户
func AutoCleanLogoutUser() {
	for {

		mUserPoolLocker.Lock()
		var nowTime = uint32(time.Now().Unix())
		for k, u := range mUserPoolMap {
			if u.IsLogin {
				continue
			}
			if u.LogoutTime > nowTime - 5 * 60 {
				continue
			}

			delete(mUserPoolMap, k)
		}
		mUserPoolLocker.Unlock()

		<-time.After(1 * time.Minute)
	}
}