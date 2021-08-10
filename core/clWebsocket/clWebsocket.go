package clWebsocket

import (
	"github.com/xiaolan580230/clws-framework/core/clDebug"
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
	clDebug.Debug("UpgradeError: status:%v reason:%v", status, reason)
}


// 启动RPC服务
//@author xiaolan
//@param _port 端口
func Serve(_path string, _port uint32) error {

	// 启动自动清理线程
	go clUserPool.AutoCleanLogoutUser()

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
				clDebug.Err("ReadMessage Error: %v", err)
			} else {
				clDebug.Info("用户断开连线...")
			}
			clUserPool.RemoveUser(uInfo.ConnId)
			break
		}

		clDebug.Info("收到消息: %v", string(buffer))
		// 心跳回应
		if msgType == websocket.PingMessage {
			_ws.WriteMessage(websocket.PongMessage, []byte{})
			continue
		}

		var requestObj clPacket.ClPacketReq
		var unMarshaErr = json.Unmarshal(buffer, &requestObj)
		if unMarshaErr != nil {
			clDebug.Err("反序列化请求失败! 错误:%v 内容:(%v)", unMarshaErr, string(buffer))
			continue
		}

		ruleInfo := clRouter.GetRule(requestObj.AC)
		if ruleInfo == nil {
			clDebug.Err("找不到路由规则: %v", requestObj.AC)
			continue
		}

		isPass, params := ruleInfo.CheckParam(requestObj.Param)
		if !isPass {
			clDebug.Err("参数:%v列表检验不通过!", requestObj.Param)
			continue
		}

		if ruleInfo.Login && !uInfo.IsLogin {
			clDebug.Err("用户:%v 未登录! 无法访问需要登录的接口:%v", uInfo.ConnId, ruleInfo.Ac)
			clRouter.SendMessage(uInfo,  requestObj.SYN, "needLogin", "您还未登录!", nil)
			continue
		}

		if ruleInfo.Callback == nil {
			clDebug.Err("接口:%v 回调函数不存在!", ruleInfo.Ac)
			continue
		}

		clDebug.Debug("收到参数列表: %+v", params)
		var resp = ruleInfo.Callback(uInfo, params)
		if resp != nil {
			serverResp := clPacket.NewPacketResp(requestObj.SYN, resp)
			_ws.WriteMessage(websocket.TextMessage, []byte(serverResp))
		}
	}

	if _ws != nil {
		_ws.Close()
	}

}

