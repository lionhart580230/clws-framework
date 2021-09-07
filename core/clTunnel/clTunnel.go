package clTunnel

import (
	"fmt"
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clws-framework/core/clCommon"
	"sync"
	"time"
)

// 隧道管理器，通过隧道可以轻松实现事件分发

var mTunnelMap map[string] map[string]ClTunnelInfo
var mTunnelLocker sync.RWMutex

type ClTunnelInfo struct {
	Tunnel chan string
	Key string
	AutoDelete bool
}

func init() {
	mTunnelMap = make(map[string] map[string]ClTunnelInfo)
}



// 创建key
func createKey() string {
	return fmt.Sprintf("%v%v", uint64(time.Now().UnixNano()), uint32(clCommon.RandInt(0,0xFFFFFFFF)))
}


// 创建管道
func Create(_name string, _autoDelete bool) *ClTunnelInfo {
	mTunnelLocker.Lock()
	_, exist := mTunnelMap[_name]
	if !exist {
		mTunnelMap[_name] = make(map[string]ClTunnelInfo, 0)
	}

	var tunnelObj = ClTunnelInfo{
		Tunnel:     make(chan string),
		AutoDelete: _autoDelete,
		Key: createKey(),
	}
	mTunnelMap[_name][tunnelObj.Key] = tunnelObj
	mTunnelLocker.Unlock()
	return &tunnelObj
}


// 等待管道数据
func (this *ClTunnelInfo) Wait(_name string) string {
	select {
		case retStr := <- this.Tunnel:
			mTunnelLocker.Lock()
			delete(mTunnelMap[_name], this.Key)
			mTunnelLocker.Unlock()
			return retStr
		case <-time.After(30 * time.Second):
			return "TIMEOUT"
	}
}


// 创建一个隧道
// 并等待隧道数据返回
func CreateAndWait(_name string, _autoDelete bool) string {
	var Obj = Create(_name, _autoDelete)
	return Obj.Wait(_name)
}


func Boardcast(_name string, _data string) {
	mTunnelLocker.Lock()
	defer mTunnelLocker.Unlock()
	_, exist := mTunnelMap[_name]
	if !exist {
		clLog.Error("tunnel: %v 不存在! 无法广播数据!", _name)
		return
	}
	for k, tunnel := range mTunnelMap[_name] {
		tunnel.Tunnel <- _data
		if tunnel.AutoDelete {
			delete(mTunnelMap[_name], k)
		}
	}
}