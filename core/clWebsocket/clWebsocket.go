package clWebsocket

import (
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clws-framework/core/clPacket"
	"github.com/xiaolan580230/clws-framework/core/clRouter"
	"github.com/xiaolan580230/clws-framework/core/clUserPool"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

// 协议升级
var upgrader = websocket.Upgrader{
	Error:       UpgradeError,
	CheckOrigin: CheckOrigin,
}

// 跨域验证
func CheckOrigin(r *http.Request) bool {
	return true
}

// 升级协议错误回调
func UpgradeError(w http.ResponseWriter, r *http.Request, status int, reason error) {
	clLog.Debug("UpgradeError: status:%v reason:%v", status, reason)
}


// 写入管道对象
var mWriteChannel chan WriteObj



// 开启写入线程
func StartWriteChannel() {
	mWriteChannel = make(chan WriteObj)
	for {
		writeBuffer := <-mWriteChannel
		clUserInfo := clUserPool.GetUserById(writeBuffer.connId)
		if clUserInfo == nil {
			clLog.Error("发送数据: %v 失败! 未找到用户连线Id: %v", writeBuffer.data, writeBuffer.connId)
			break
		}
		sendErr := clUserInfo.SendMsg(writeBuffer.data)
		if sendErr != nil {
			clLog.Error("发送数据: %v 失败! 错误:%v", writeBuffer.data, sendErr)
		}
	}
}



// 启动RPC服务
//@author xiaolan
//@param _port 端口
func Serve(_path string, _port uint32) error {

	// 启动自动清理线程
	go clUserPool.AutoCleanLogoutUser()

	// 启动2个写入类型管道对象
	go StartWriteChannel()
	go StartWriteChannel()

	// websocket 服务
	http.HandleFunc("/" + _path, doWork)

	// 启动服务
	if err := http.ListenAndServe(fmt.Sprintf(":%v", _port), nil); err != nil {
		return err
	}

	return nil
}



//处理主要业务逻辑
//@author xiaolan
//@param w 返回消息体
//@param r 接收消息体
func doWork(w http.ResponseWriter, r *http.Request) {

	// 跨域支持
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// 协议升级
	_ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade error: %v\n", err)
		return
	}

	uInfo := clUserPool.AddNewUser(_ws, false)
	// 消息循环
	for {
		// 获取消息
		msgType, buffer, err := _ws.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "EOF") && strings.Contains(err.Error(), "close 1005") {
				clLog.Error("ReadMessage Error: %v", err)
			} else {
				clLog.Info("用户[%v]断开连线...", uInfo.ConnId)
			}
			clUserPool.RemoveUser(uInfo.ConnId)
			break
		}

		clLog.Debug("收到消息: %v", string(buffer))
		// 心跳回应
		if msgType == websocket.PingMessage {
			_ws.WriteMessage(websocket.PongMessage, []byte{})
			continue
		}

		var requestObj clPacket.ClPacketReq
		var unMarshaErr = json.Unmarshal(buffer, &requestObj)
		if unMarshaErr != nil {
			clLog.Error("反序列化请求失败! 错误:%v 内容:(%v)", unMarshaErr, string(buffer))
			continue
		}

		ruleInfo := clRouter.GetRule(requestObj.AC)
		if ruleInfo == nil {
			clLog.Error("找不到路由规则: %v", requestObj.AC)
			continue
		}

		isPass, params := ruleInfo.CheckParam(requestObj.Param)
		if !isPass {
			clLog.Error("参数:%v列表检验不通过!", requestObj.Param)
			clRouter.SendMessage(uInfo, "paramError", "参数错误", nil)
			continue
		}

		if ruleInfo.Login && !uInfo.IsLogin {
			clLog.Error("用户:%v 未登录! 无法访问需要登录的接口:%v", uInfo.ConnId, ruleInfo.Ac)
			clRouter.SendMessage(uInfo, "needLogin", "您还未登录", nil)
			continue
		}

		if ruleInfo.Callback == nil {
			clLog.Error("接口:%v 回调函数不存在!", ruleInfo.Ac)
			continue
		}

		clLog.Debug("收到参数列表: %+v", params)
		// 启动线程处理
		go func(_connId uint64) {
			var resp = ruleInfo.Callback(uInfo, params)
			if resp != nil {
				mWriteChannel <- WriteObj{
					data: clPacket.NewPacketResp( resp ),
					connId: _connId,
				}
			}
		} (uInfo.ConnId)

	}

	// 从用户池中移除
	clUserPool.RemoveUser(uInfo.ConnId)
}

