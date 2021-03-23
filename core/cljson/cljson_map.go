package cljson

import (
	"clws-framework/core/clCommon"
	"regexp"
	"strings"
)

// 遍历数组
func (js *JsonMap) Each(hf func(key string, val *JsonStream)) {
	for key, val := range *js {
		hf(key, val.data)
	}
}

// 获取指定的对象
func (js *JsonMap) Find(key string) *JsonStream {
	if val, Ok := (*js)[key]; Ok {
		return val.data
	}
	return nil
}


// key是否存在
func (js *JsonMap) IsSet(_key string) bool {
	_, ok := (*js)[_key]
	return ok
}


// 重置为标准库的map[string] interface{} 去掉了反斜杠
func (js * JsonMap) ToCustom2() map[string] string {
	var CustomMap = make(map[string] string)
	reg,_ := regexp.Compile(`(^\\)*\\\\(^\\)*`)
	for key, val := range *js {
		if val.data == nil {
			continue
		}
		res := reg.ReplaceAllString(val.data.ToStr(),"_replace_xiegang_") // 换掉 \\
		res = strings.ReplaceAll(res,"\\n","_replace_huanhang_") // 换掉\n
		res = strings.ReplaceAll(res,"\\","") // 换掉剩余的\
		res = strings.ReplaceAll(res,"_replace_xiegang_","\\") // 还原斜杠
		res = strings.ReplaceAll(res,"_replace_huanhang_","\\n") // 还原换行
		CustomMap[key] = res
	}
	return CustomMap
}

// 重置为标准库的map[string] interface{}
func (js * JsonMap) ToCustom() map[string] string {
	var CustomMap = make(map[string] string)
	for key, val := range *js {
		if val.data == nil {
			continue
		}
		CustomMap[key] = val.data.ToStr()
	}
	return CustomMap
}

// 重置为标准库的map[string] string
func (js * JsonMap) Tostring() map[string] string{
	var CustomMap = make(map[string] string)
	for key, val := range *js {
		CustomMap[key] = val.data.ToStr()
	}
	return CustomMap
}


// 重置为标准库的map[string] interface{}
func (js * JsonMap) ToTree() map[string] interface{}{
	//fmt.Printf(">> %v\n", js.Tostring())
	var CustomMap = make(map[string] interface{})
	for key, val := range *js {
		//fmt.Printf("key: %v => val: %v\n",key, val)  ==> Log 太多太难受， 到外面打印去 Atin 05-14
		CustomMap[key] = val.StackParseTree()
	}
	return CustomMap
}

// 递归查找树结构
func (this *jsonItem) StackParseTree() interface{} {

	//fmt.Printf(">> STR: %v => TYPE: %v\n", string(this.data.data), this.data.dataType)

	if this.data.dataType == JSON_TYPE_MAP {
		this_map := this.data.ToMap()
		var customMap = make(map[string] interface{})
		for key, val := range *this_map {
			customMap[key] = val.StackParseTree()
		}
		return customMap
	} else if this.data.dataType == JSON_TYPE_ARR {
		this_arr := this.data.ToArray()
		var customArr = make([] interface{}, 0)
		for _, val := range *this_arr {
			customArr= append(customArr, val.StackParseTree() )
		}
		return customArr
	} else if this.data.dataType == JSON_TYPE_BOOL {
		return clCommon.Bool(this.data.ToStr())
	} else if this.data.dataType == JSON_TYPE_INT {
		return clCommon.Float64(this.data.ToStr(), 0)
	} else if this.data.dataType == JSON_TYPE_NULL {
		return nil
	}
	return this.data.ToStr()
}


// 获取指定下标并转换成string类型
// @param key string 下标
// @param def string 默认值
// @return string 返回指定的之
func (js * JsonMap) GetStr(key string, def string) string {
	if val, ok := (*js)[key]; ok {
		return val.data.ToStr()
	}
	return ""
}

// 获取指定下标并转换成前后去空的string类型
// @param key string 下标
// @param def string 默认值
// @return string 返回指定的之
func (js * JsonMap) GetStrTrim(key string, def string) string {
	if val, ok := (*js)[key]; ok {
		return strings.TrimSpace(val.data.ToStr())
	}
	return ""
}

// 获取指定下标并转换成uint32类型
func (js * JsonMap) GetUint32(key string, def uint32) uint32 {
	if val, ok := (*js)[key]; ok {
		return clCommon.Uint32(val.data.ToStr(), 0)
	}
	return def
}

// 获取指定下标并转换成uint32类型
func (js * JsonMap) GetUint64(key string, def uint64) uint64 {
	if val, ok := (*js)[key]; ok {
		return clCommon.Uint64(val.data.ToStr(), 0)
	}
	return def
}

// 获取指定下标并转换成int32类型
func (js * JsonMap) GetInt32(key string, def int32) int32 {
	if val, ok := (*js)[key]; ok {
		return clCommon.Int32(val.data.ToStr(), 0)
	}
	return def
}

// 获取指定下标并转换成float32类型
func (js * JsonMap) GetFloat32(key string, def float32) float32 {
	if val, ok := (*js)[key]; ok {
		return clCommon.Float32(val.data.ToStr(), 0)
	}
	return def
}


// 获取指定下标并转换成float64类型
func (js * JsonMap) GetFloat64(key string, def float64) float64 {
	if val, ok := (*js)[key]; ok {
		return clCommon.Float64(val.data.ToStr(), 0)
	}
	return def
}


// 获取指定下标并转换成bool类型
func (js * JsonMap) GetBool(key string, def bool) bool {
	if val, ok := (*js)[key]; ok {
		return clCommon.Bool(val.data.ToStr())
	}
	return def
}

func (js * JsonMap) DelKey(key string) *JsonMap {
	if _, ok := (*js)[key]; ok {
		delete(*js, key)
		return js
	}
	return nil
}


func (js * JsonMap) IsEmpty() bool {
	if len(*js) == 0 {
		return true
	}
	return false
}