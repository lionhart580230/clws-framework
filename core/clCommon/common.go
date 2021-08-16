package clCommon

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 强制转换为int32类型
func Int32(_val string, _def int32) int32 {
	i, e := strconv.ParseInt(_val, 10, 32)
	if e != nil {
		return _def
	}
	return int32(i)
}


// 强制转换为int32类型
func Int(_val string, _def int) int {
	i, e := strconv.ParseInt(_val, 10, 32)
	if e != nil {
		return _def
	}
	return int(i)
}


// 强制转换为int64类型
func Int64(_val string, _def int64) int64 {
	i, e := strconv.ParseInt(_val, 10, 64)
	if e != nil {
		return _def
	}
	return int64(i)
}


// 强制转换为uint32类型
func Uint32(_val string, _def uint32) uint32 {
	i, e := strconv.ParseUint(_val, 10, 32)
	if e != nil {
		return _def
	}
	return uint32(i)
}


// 强制转换为uint64类型
func Uint64(_val string, _def uint64) uint64 {
	i, e := strconv.ParseUint(_val, 10, 64)
	if e != nil {
		return _def
	}
	return uint64(i)
}


// 强制转换为float32类型
func Float32(_val string, _def float32) float32 {
	f, e := strconv.ParseFloat(_val, 32)
	if e != nil {
		return _def
	}
	return float32(f)
}


// 强制转换为float64类型
func Float64(_val string, _def float64) float64 {
	f, e := strconv.ParseFloat(_val, 64)
	if e != nil {
		return _def
	}
	return f
}


// 强制转换为Bool类型
func Bool(_val string) bool {
	var val = strings.ToUpper(_val)
	if val == "TRUE" || val == "ON" || val == "YES" || val == "1" {
		return true
	}
	return false
}


func HtmlSpecialChars(val string) string {

	r, _ := regexp.Compile(`[\&]`)
	val = string(r.ReplaceAll([]byte(val), []byte("&amp;")))
	r, _ = regexp.Compile(`[\>]`)
	val = string(r.ReplaceAll([]byte(val), []byte("&gt;")))
	r, _ = regexp.Compile(`[\<]`)
	val = string(r.ReplaceAll([]byte(val), []byte("&lt;")))
	r, _ = regexp.Compile(`[\"]`)
	val = string(r.ReplaceAll([]byte(val), []byte("&quot;")))
	r, _ = regexp.Compile(`[\']`)
	val = string(r.ReplaceAll([]byte(val), []byte("&#039;")))
	return val
}


func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 取一个随机数
func RandInt(min int, max int) int {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	var nan = int(0)
	for i := 0; i < 10; i++ {
		nan += int(rand.Int63())
		nan = (nan % (max - min)) + (min)
	}
	return nan
}