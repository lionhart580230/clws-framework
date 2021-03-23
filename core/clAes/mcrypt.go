package clAes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"net/url"
)

// 服务端解密程序
// 使用 AES 加密方式

type SkyPacket struct {
	Iv string `json:"iv"`
	Value string `json:"value"`
}


//@author xiaolan
//@lastUpdate 2019-09-29
//@comment 解密程序
func Decode(_buffer []byte, _key []byte) []byte {
	dData := base64Decode([]byte(_buffer))
	var Packet SkyPacket

	bufferResp := make([]byte, 0)
	// 解析失败
	err := json.Unmarshal(dData, &Packet)
	if err != nil {
		return nil
	}

	dIv := base64Decode( []byte(Packet.Iv) )
	dValue := base64Decode( []byte(Packet.Value) )

	_cipher, err := aes.NewCipher([]byte(_key))
	if err != nil {
		return nil
	}

	blockMode := cipher.NewCBCDecrypter(_cipher, []byte(dIv))

	origData := make([]byte, len(dValue))
	blockMode.CryptBlocks(origData, dValue)
	// 解析
	bufferResp = PKCS5UnPadding(origData)

	return bufferResp
}


//@author xiaolan
//@lastUpdate 2019-09-29
//@comment 加密程序
func Encode(_buffer []byte, _key []byte) []byte {

	_cipher, err := aes.NewCipher([]byte(_key))
	if err != nil {
		return nil
	}

	iv := RandomBlock()
	value := _buffer

	// 生成CBC加密对象
	blockMode := cipher.NewCBCEncrypter(_cipher, iv)

	// 填充字节
	dValue := PKCS5Padding(value)
	origData := make([]byte, len(dValue))

	// 加密
	blockMode.CryptBlocks(origData, dValue)

	// 组装结构
	var skyPacket = SkyPacket{
		Iv:    string(base64Encode(iv)),
		Value: string(base64Encode(origData)),
	}

	// 生成json字串
	finallyData, err := json.Marshal(skyPacket)
	if err != nil {
		fmt.Printf(">> 生成json字符串失败! 错误: %v\n", err)
		return nil
	}

	// base64加密
	return []byte(base64Encode(finallyData))
}



// 解密数据
func DecodeParam(_params string, _key string) url.Values {
	paramStr := string(Decode([]byte(_params), []byte(_key)))

	mValues, err := url.ParseQuery(paramStr)
	if err != nil {
		fmt.Printf(">> ParseQuery Error: %v\n", err)
		return nil
	}
	return mValues
}