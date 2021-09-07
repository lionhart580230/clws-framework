package clPacket

import (
	"encoding/json"
	"time"
)

// 请求结构体
type ClPacketReq struct {
	AC string `json:"ac"`			// 请求路由名称
	Timestamp uint32 `json:"ts"`	// 请求发起时间戳（s级）
	Header string `json:"h"`		// 请求头部
	Param string `json:"p"`			// 请求附带参数列表 UrlEncode
	Sign string `json:"sg"`			// 签名串
}


// 回应结构体
type ClPacketResp struct {
	RP string `json:"rp"`			// 响应路由名称
	TimeStamp uint32 `json:"ts"`	// 时间戳
	Param interface{} `json:"p"`	// 响应内容 json的base64加密串
	Tips string `json:"tip"`		// 提示
}


type RuleCBResp struct {
	RC string
	Param string
	Data interface{}
	RoomId uint32
}


// 生成服务器响应包
func NewPacketResp(_data *RuleCBResp) string {

	var obj = ClPacketResp{
		RP:        _data.RC,
		TimeStamp: uint32(time.Now().Unix()),
		Param:     _data.Data,
		Tips: 	   _data.Param,
	}
	var packetStr []byte
	packetStr, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(packetStr)
}

