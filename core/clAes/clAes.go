package clAes

import (
	"bytes"
	"encoding/base64"
	"time"
	"math/rand"
)

// 生成16位随机字符串
func RandomBlock() []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/"
	bytes := []byte(str)
	result := []byte{}
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Int31n(int32(len(str)))])
	}
	return result
}


func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(_origData []byte) []byte {
	padnum := 16 - len(_origData) % 16
	pad := bytes.Repeat([]byte{ byte(padnum) }, padnum)
	return append(_origData, pad...)
}

// Base64解密
func base64Decode(data []byte) []byte {

	lenOfData := len(data)
	if lenOfData%4 > 0{
		data = append(data, bytes.Repeat([]byte("="), 4-(lenOfData%4))...)
	}
	srcByte, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil{
		return []byte{}
	}

	return srcByte
}


// Base64加密
func base64Encode(data []byte) []byte {

	srcByte := base64.StdEncoding.EncodeToString(data)

	return []byte(srcByte)
}

