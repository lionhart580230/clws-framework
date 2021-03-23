package clPacket

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// 请求结构体
type ClPacketReq struct {
	AC string `json:"ac"`			// 请求路由名称
	Timestamp uint32 `json:"ts"`	// 请求发起时间戳（s级）
	Param string `json:"p"`			// 请求附带参数列表 UrlEncode
	Sign string `json:"sg"`			// 签名串
	SYN uint32 `json:"syn"`			// 请求ID, 由前台生成
}


// 回应结构体
type ClPacketResp struct {
	RP string `json:"rp"`			// 响应路由名称
	TimeStamp uint32 `json:"ts"`	// 时间戳
	Param string `json:"p"`			// 响应内容 json的base64加密串
	Tips string `json:"tip"`		// 提示
	Sign string `json:"sg"`			// 签名串
	ACK uint32 `json:"ack"`			// 响应ID, 与请求对应, 推送消息为0
}

type RuleCBResp struct {
	RC string
	Param string
	Data interface{}
}


// 生成服务器响应包
func NewPacketResp(_ack uint32, _data *RuleCBResp) string {
	var paramStr = ""
	if _data.Data != nil {
		jsonBytes, err := json.Marshal(_data.Data)
		if err != nil {
			return ""
		}
		paramStr = base64.StdEncoding.EncodeToString(jsonBytes)
	}

	var obj = ClPacketResp{
		RP:        _data.RC,
		TimeStamp: uint32(time.Now().Unix()),
		Param:     paramStr,
		Sign:      "",
		ACK:       _ack,
		Tips: 	   _data.Param,
	}
	var packetStr []byte
	packetStr, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(packetStr)
}

