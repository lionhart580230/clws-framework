package clmysql

import (
	"github.com/xiaolan580230/clws-framework/core/clCommon"
	"strings"
)

// 回传数据类型解析

// 获取Int32类型的值
func (res *TdbResult) GetInt32(key string, def int32) int32 {
	if val, find := (*res)[key]; find {
		return clCommon.Int32(val, 0)
	}
	return def
}

// 获取Uint32
func (res *TdbResult) GetUint32(key string, def uint32) uint32 {
	if val, find := (*res)[key]; find {
		return clCommon.Uint32(val, 0)
	}
	return def
}

// 获取 Uint64
func (res *TdbResult) GetUint64(key string, def uint64) uint64 {
	if val, find := (*res)[key]; find {
		return clCommon.Uint64(val, 0)
	}
	return def
}

// 获取 Uint64
func (res *TdbResult) GetInt64(key string, def int64) int64 {
	if val, find := (*res)[key]; find {
		return clCommon.Int64(val, 0)
	}
	return def
}

// 获取Float32
func (res *TdbResult) GetFloat32(key string, def float32) float32 {
	if val, find := (*res)[key]; find {
		return clCommon.Float32(val, 0)
	}
	return def
}

// 获取Float64
func (res *TdbResult) GetFloat64(key string, def float64) float64 {
	if val, find := (*res)[key]; find {
		return clCommon.Float64(val, 0)
	}
	return def
}

// 获取Bool
func (res *TdbResult) GetBool(key string, def bool) bool {
	if val, find := (*res)[key]; find {
		if strings.ToLower(val) == "true" || strings.ToLower(val) == "yes" || strings.ToLower(val) == "on" || val == "1" {
			return true
		}
	}
	return false
}

// 获取字符串
func (res *TdbResult) GetStr(key string, def string) string {
	if val, find := (*res)[key]; find {
		return val
	}
	return def
}

/*
	循环 TdbResult
 */
func (res *TdbResult) Each(hf func(key string, value string) bool) {
	for k, val := range *res {
		if ok := hf(k, val); !ok {
			return
		}
	}
}


//// 获取日期格式
//func (res *TdbResult) GetDate(key string, def string) string {
//
//	if val, find := (*res)[key]; find {
//		if common.Uint32(val) == 0 {
//			return def
//		}
//		return cltime.GetDateByFormat(common.Uint32(val), "2006-01-02")
//	} else {
//		return def
//	}
//
//}
//
//// 获取时间格式
//func (res *TdbResult) Gettime(key string, def string) string {
//
//	if val, find := (*res)[key]; find {
//		if common.Uint32(val) == 0 {
//			return def
//		}
//		return cltime.GetDateByFormat(common.Uint32(val), "15:04:05")
//	} else {
//		return def
//	}
//
//}
//
//// 获取日期时间格式
//func (res *TdbResult) GetDatetime(key string, def string) string {
//
//	if val, find := (*res)[key]; find {
//		if common.Uint32(val) == 0 {
//			return def
//		}
//		return cltime.GetDateByFormat(common.Uint32(val), "2006-01-02 15:04:05")
//	} else {
//		return def
//	}
//
//}