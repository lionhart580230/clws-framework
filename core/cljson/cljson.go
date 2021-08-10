package cljson

import (
	"bytes"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xiaolan580230/clws-framework/core/clCommon"
	"reflect"
	"regexp"
	"strings"
)

type M map[string] interface{}
type A []interface{}

// 创建一个新对象
// {"":""}
func New(bs []byte) (*JsonStream) {
	var js JsonStream
	js.data = bs
	js.data = bytes.Replace(js.data, []byte{10}, []byte{}, -1)		// 去除tab
	js.data = bytes.Replace(js.data, []byte{13}, []byte{}, -1)		// 去除tab
	js.dataLength = uint32(len(js.data))

	if !js.IsValidJson() {
		return nil
	}

	return &js
}


// 创建一个新的实体
func CreateBy(v interface{}) (*JsonStream) {

	var js JsonStream
	var err error
	js.data, err = jsoniter.Marshal(v)
	if err != nil {
		fmt.Printf(">> 创建json失败! marshal错误: %v", err)
		return &js
	}
	temp := js.data
	js.dataLength = uint32(cap(temp))
	if !js.IsValidJson() {
		fmt.Printf(">> 创建JSON失败! (%v)\n", string(temp))
		return &js
	}
	return &js
}



// 检测某个key是否存在
func (js *JsonStream) IsSet(_key string) bool {
	if js.dataMap == nil {
		return false
	}

	_, exists := js.dataMap[_key]
	return exists
}



// 判断是否是json_array
func (js *JsonStream) IsValidArray() bool {

	lessArray := 1
	lessObject := 0
	isInString := false
	cacheArray := make([]string, 0)
	itemBytes := make([]byte, 0)
	newBytes := bytes.Buffer{}

	valType := JSON_TYPE_NIL

	if len(js.data) == 2 {
		return true
	}

	for i, v := range js.data {
		if i == 0 {
			newBytes.WriteByte(v)
			continue
		}

		// 切换字符串模式
		if v == '"' && (i <= 0 || js.data[i-1] != '\\') {
			isInString = !isInString
			if valType == JSON_TYPE_NIL {
				valType = JSON_TYPE_STR
			}
			itemBytes = append(itemBytes, []byte("\"")...)
			newBytes.WriteByte(v)
			continue
		}

		isPass := false
		// 字符串处理
		if isInString {
			itemBytes = append(itemBytes, v)
			newBytes.WriteByte(v)
			continue
		}

		if v != ' ' {
			newBytes.WriteByte(v)
		}

		// 各种结构解析
		switch v {
		case '{':
			if valType == JSON_TYPE_NIL {
				if lessObject == 1 && lessArray == 0 {
					valType = JSON_TYPE_MAP
				}
			}
			lessObject ++
		case '[':
			if valType == JSON_TYPE_NIL {
				if lessObject == 1 && lessArray == 0 {
					valType = JSON_TYPE_ARR
				}
			}
			lessArray ++
		case '}':
			lessObject --
			if lessObject < 0 {
				return false
			}

		case ']':
			lessArray --

			if lessArray == 0 {
				vtype := JSON_TYPE_NIL
				vtype, itemBytes = checkType(itemBytes)
				if vtype == JSON_TYPE_NIL {
					return false
				}

				cacheArray = append(cacheArray, string(itemBytes))
				itemBytes = make([]byte, 0)
				isPass = true
				valType = JSON_TYPE_NIL
			}
		}

		if !isPass {
			if v == ',' {

				if lessObject > 0 || lessArray > 1 {
					itemBytes = append(itemBytes, v)
					continue
				}
				vtype := JSON_TYPE_NIL
				vtype, itemBytes = checkType(itemBytes)
				if vtype == JSON_TYPE_NIL {

					return false
				}

				cacheArray = append(cacheArray, string(itemBytes))
				itemBytes = make([]byte, 0)
				valType = JSON_TYPE_NIL

			} else {
				if valType != JSON_TYPE_NIL {
					itemBytes = append(itemBytes, v)
				} else if v != ' ' {
					itemBytes = append(itemBytes, v)
				}
			}
		}
	}


	if lessArray != 0 || lessObject != 0 || isInString {
		return false
	}

	js.dataArray = cacheArray
	js.data = newBytes.Bytes()
	return true
}



// 判断是否是json_map
func (js *JsonStream) IsValidMap() bool {

	lessArray := 0
	lessObject := 1
	isInString := false
	lastKey := ""
	lastVal := make([]byte, 0)
	newBytes := bytes.Buffer{}
	keyMode := true
	valType := JSON_TYPE_NIL

	if len(js.data) == 2 {
		return true
	}

	cacheMap := make(map[string] string)
	for i, v := range js.data {
		if i == 0 {
			newBytes.WriteByte(v)
			continue
		}

		// 切换字符串模式
		if v == '"' && (i <= 0 || js.data[i-1] != '\\') {
			isInString = !isInString
			if !keyMode {
				if valType == JSON_TYPE_NIL {
					valType = JSON_TYPE_STR
				}
				lastVal = append(lastVal, []byte("\"")...)
			}
			newBytes.WriteByte(v)
			continue
		}

		isPass := false
		// 字符串处理
		if isInString {
			if keyMode {		// 查找键模式
				if v == ' ' {
					return false
				}
				lastKey += string(v)
			} else {			// 查找值模式
				lastVal = append(lastVal, v)
			}
			newBytes.WriteByte(v)
			continue
		}

		if v != ' ' {
			newBytes.WriteByte(v)
		}

		// 各种结构解析
		switch v {
			case '{':
				if valType == JSON_TYPE_NIL {
					if lessObject == 1 && lessArray == 0 {
						valType = JSON_TYPE_MAP
					}
				}
				lessObject ++
			case '[':
				if valType == JSON_TYPE_NIL {
					if lessObject == 1 && lessArray == 0 {
						valType = JSON_TYPE_ARR
					}
				}
				lessArray ++
			case '}':
				lessObject --
				if lessObject <= 0 && i != int(js.dataLength)-1 {
					return false
				}

				if lessObject == 0 {
					vtype := JSON_TYPE_NIL
					vtype, lastVal = checkType([]byte(lastVal))
					if vtype == JSON_TYPE_NIL {

						return false
					}

					cacheMap[lastKey] = string(lastVal)
					lastKey = ""
					lastVal = make([]byte, 0)
					keyMode = true
					isPass = true
					valType = JSON_TYPE_NIL
				}

			case ']':
				lessArray --
			case ':':
				if lessObject == 1 && lessArray == 0 {
					keyMode = false
					isPass = true
				}
			}


		if !isPass {

			if keyMode {
				if v != ' ' {
					return false
				}

			} else {

				if v == ',' {
					if lessObject > 1 || lessArray > 0 {
						lastVal = append(lastVal, v)
						continue
					}

					if lastKey == "" || len(lastVal) == 0 {
						return false
					}

					vtype := JSON_TYPE_NIL
					vtype, lastVal = checkType([]byte(lastVal))
					if vtype == JSON_TYPE_NIL {
						return false
					}

					cacheMap[lastKey] = string(lastVal)
					lastKey = ""
					lastVal = []byte{}
					keyMode = true
					valType = JSON_TYPE_NIL
				} else {
					if valType != JSON_TYPE_NIL {
						lastVal = append(lastVal, v)
					} else if v != ' ' {
						lastVal = append(lastVal, v)
					}
				}
			}
		}
	}

	if lessArray != 0 || lessObject != 0 || isInString {
		return false
	}

	js.dataMap = cacheMap
	//for key, val := range cacheMap {
	//	fmt.Printf(">> %v => %v\n", key, val)
	//}
	js.data = newBytes.Bytes()
	return true
}


// 检查类型
func checkType(data []byte) (int, []byte) {
	if len(data) < 1 {
		return JSON_TYPE_NIL, nil
	}
	//fmt.Printf(">> type: (%v)\n", string(data))

	data = bytes.TrimSpace(data)
	end := len(data)-1
	switch data[0] {
	case '{':
		if data[end] == '}' {
			return JSON_TYPE_MAP, data
		}
	case '[':
		if data[end] == ']' {
			return JSON_TYPE_ARR, data
		}
	case '"':
		if data[end] == '"' {
			return JSON_TYPE_STR, data
		}
	default:
		reg, _ := regexp.Compile(`^(\-)?[0-9]+(\.[0-9]+)?(e[0-9]+)?$`)
		reg2, _ := regexp.Compile(`^(true|false|TRUE|FALSE)$`)
		regE, _ := regexp.Compile(`^(\-)?[0-9](\.[0-9]+)E(\-)?[0-9]$`)

		if reg2.Match(data) {
			return JSON_TYPE_BOOL, data
		}

		if reg.Match(data) || regE.Match(data) {
			return JSON_TYPE_INT, data
		}

		m, _ := regexp.Match(`^(null|NULL)$`, data)
		if m {
			return JSON_TYPE_NULL, data
		}
	}
	return JSON_TYPE_NIL, nil
}


// 判断json结构是否有效
func (js *JsonStream) IsValidJson() bool {
	if js == nil {
		return false
	}

	if js.dataLength == 0 {

		return false
	}

	end := js.data[js.dataLength-1]

	if js.data[0] == '[' && end == ']' {
		js.dataType = JSON_TYPE_ARR
		return js.IsValidArray()
	} else if js.data[0] == '{' && end == '}' {
		js.dataType = JSON_TYPE_MAP
		return js.IsValidMap()
	} else if js.data[0] == '"' && end == '"' {
		js.dataType = JSON_TYPE_STR
		js.data = js.data[1:js.dataLength-1]
		js.dataLength = uint32(len(js.data))
		return true
	} else if string(js.data[:2]) == `\"` && string(js.data[js.dataLength-2:]) == `\"` {
		js.dataType = JSON_TYPE_STR
		js.data = js.data[2:js.dataLength-2]
		js.dataLength = uint32(len(js.data))
		return true
	}

	// 判断是否是数值
	reg, _ := regexp.Compile(`^(\-)?[0-9]+(\.[0-9]+)?(e[0-9]+)?$`)		// 整数/小数
	reg2, _ := regexp.Compile(`^(true|false|TRUE|FALSE)$`)	// 布尔
	regE, _ := regexp.Compile(`^(\-)?[0-9](\.[0-9]+)E(\-)?[0-9]$`)

	if reg2.Match(js.data) {
		js.dataType = JSON_TYPE_BOOL
		return true
	}

	if reg.Match(js.data) || regE.Match(js.data) {
		js.dataType = JSON_TYPE_INT
		return true
	}

	m, _ := regexp.Match(`^(null|NULL)$`, js.data)
	if m {
		js.dataType = JSON_TYPE_NULL
		return true
	}

	return false
}


// 转换成字符串
func (js *JsonStream) ToStr() string {
	var min = js.lastLockBegin
	var max = js.lastLockEnd

	if max == 0 {
		max = len(js.data)
	}
	if max >= 1 {
		if js.data[min] == '"' && js.data[max-1] == '"' {
			min = min + 1
			max = max - 1
		}
	}
	return string(js.data[min:max])
}


// 转换成字符串
func (js *JsonStream) ToFloat32() float32 {
	return clCommon.Float32(js.ToStr(), 0)
}

// 转换成字符串
func (js *JsonStream) ToFloat64() float64 {
	return clCommon.Float64(js.ToStr(), 0)
}

// 转换成字符串
func (js *JsonStream) ToInt32() int32 {
	return clCommon.Int32(js.ToStr(), 0)
}

// 转换成字符串
func (js *JsonStream) ToUint32() uint32 {
	return clCommon.Uint32(js.ToStr(), 0)
}

// 转换成字符串
func (js *JsonStream) ToInt64() int64 {
	return clCommon.Int64(js.ToStr(), 0)
}

// 转换成字符串
func (js *JsonStream) ToUint64() uint64 {
	return clCommon.Uint64(js.ToStr(), 0)
}

// 转换成bool
func (js *JsonStream) ToBool() bool {
	return clCommon.Bool(js.ToStr())
}

// 获取字节集
func (js *JsonStream) ToJSON() []byte {
	var begin = js.lastLockBegin
	var end = js.lastLockEnd
	if end == 0 {
		end = len(js.data)
	}
	return js.data[begin:end]
}

// 变成Arr
func (js *JsonStream) ToArray()  *JsonArray {

	var jsonArr = make(JsonArray, 0)
	for _, val := range js.dataArray {
		jsonArr = append(jsonArr, jsonItem{ New([]byte(val)) })
	}
	return &jsonArr
}

// 变成map
func (js *JsonStream) ToMap() *JsonMap {

	var jsonMap = make(JsonMap, 0)
	for mKey, mVal := range js.dataMap {
		jsonMap[mKey] = jsonItem { New([]byte(mVal)) }
	}

	return &jsonMap
}


// 获取json中的字符串
func (js *JsonStream) GetStr(key string, subkey... string) string {

	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToStr()
		}
		return ""
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return ""
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return ""
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return ""
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return ""
	}
	return finaly.ToStr()

	//var findBytes = make([]byte, 0)
	//if js.lastLockEnd == 0 {
	//	findBytes = findLastValue(js.data, []byte(key))
	//} else {
	//	findBytes = findLastValue(js.data[js.lastLockBegin:js.lastLockEnd], []byte(key))
	//}
	//
	//if findBytes == nil {
	//	return ""
	//}
	//
	//for _, val := range subkey {
	//	if val == "" {
	//		continue
	//	}
	//	findBytes = findLastValue(findBytes, []byte(val))
	//	if findBytes == nil {
	//		return ""
	//	}
	//}
	//return string(bytes.Replace(findBytes, []byte(`\`), []byte{}, -1))
}

// 获取json中的小数
func (js *JsonStream) GetFloat32(key string, subkey... string) float32 {

	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToFloat32()
		}
		return 0
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return 0
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return 0
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return 0
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return 0
	}
	return finaly.ToFloat32()

	//
	//var findBytes = make([]byte, 0)
	//
	//if js.lastLockEnd == 0 {
	//	findBytes = findLastValue(js.data, []byte(key))
	//} else {
	//	findBytes = findLastValue(js.data[js.lastLockBegin:js.lastLockEnd], []byte(key))
	//}
	//
	//if findBytes == nil {
	//	return 0
	//}
	//
	//for _, val := range subkey {
	//	findBytes = findLastValue(findBytes, []byte(val))
	//	if findBytes == nil {
	//		return 0
	//	}
	//}
	//
	//return common.Float32(string(findBytes))
}

// 获取json中的整数
func (js *JsonStream) GetInt32(key string, subkey... string) int32 {

	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToInt32()
		}
		return 0
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return 0
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return 0
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return 0
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return 0
	}
	return finaly.ToInt32()
}

// 获取json中的整数
func (js *JsonStream) GetInt64(key string, subkey... string) int64 {

	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToInt64()
		}
		return 0
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return 0
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return 0
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return 0
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return 0
	}
	return finaly.ToInt64()
}

// 获取json中的整数
func (js *JsonStream) GetUint64(key string, subkey... string) uint64 {
	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToUint64()
		}
		return 0
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return 0
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return 0
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return 0
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return 0
	}
	return finaly.ToUint64()
}

// 获取json中的整数
func (js *JsonStream) GetUint32(key string, subkey... string) uint32 {
	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToUint32()
		}
		return 0
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return 0
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return 0
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return 0
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return 0
	}
	return finaly.ToUint32()
}

// 获取json中的小数
func (js *JsonStream) GetFloat64(key string, subkey... string) float64 {
	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToFloat64()
		}
		return 0
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return 0
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return 0
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return 0
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return 0
	}
	return finaly.ToFloat64()
}

// 获取json中的布尔类型
func (js *JsonStream) GetBool(key string, subkey... string) bool {
	val := js.dataMap[key]
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToBool()
		}
		return false
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return false
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return false
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return false
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return false
	}
	return finaly.ToBool()
}

// 获取json中的数组
func (js *JsonStream) GetArray(key string, subkey... string) *JsonArray {
	var findBytes = make([]byte, 0)
	if js.lastLockEnd == 0 {
		findBytes = findLastValue(js.data, []byte(key))
	} else {
		findBytes = findLastValue(js.data[js.lastLockBegin:js.lastLockEnd], []byte(key))
	}
	if findBytes == nil {
		return nil
	}

	for _, val := range subkey {
		findBytes = findLastValue(findBytes, []byte(val))
		if findBytes == nil {
			return nil
		}
	}

	var arrItem = parseArray(findBytes)
	var jsonArr = make(JsonArray, 0)

	for _, val := range arrItem {
		jsonArr = append(jsonArr, jsonItem{ New(val) })
	}
	return &jsonArr
}

// 获取json中的Map
func (js *JsonStream) GetMap(key string, subkey... string) *JsonMap {

	val := js.dataMap[key]
	var jsonMap = make(JsonMap)
	if len(subkey) == 0 {
		item := New([]byte(val))
		if item != nil {
			return item.ToMap()
		}
		return &jsonMap
	}

	cacheVal := val
	for _, key := range subkey {
		item := New([]byte(cacheVal))
		if item == nil {
			return &jsonMap
		}
		itemMap := item.ToMap()
		if itemMap == nil  {
			return &jsonMap
		}

		cacheVal := itemMap.GetStr(key, "")
		if cacheVal == "" {
			return &jsonMap
		}
	}

	finaly := New([]byte(cacheVal))
	if finaly == nil {
		return &jsonMap
	}
	return finaly.ToMap()
}


// 搜索指定key的位置
func findKeyPos(target []byte, key []byte) int {
	targetLen := len(target)
	keyLen := len(key)
	if keyLen >= targetLen {
		return -1
	}

	for i:=1; i<targetLen-keyLen; i++ {
		if bytes.Equal(target[i:keyLen+i], key) && target[i-1] == '"' && target[keyLen+i] == '"' {
			return i-1
		}
	}
	return -1
}


// 解析数组
func parseArray(target []byte) [][]byte{
	lenTarget := len(target)
	if lenTarget <= 2 {
		return make([][]byte, 0)
	}
	if target[0] != '[' || target[lenTarget-1] != ']' {
		return make([][]byte, 0)
	}
	isInString := false	// 是否在引号内部
	leftCount := 1
	findPos := 1		// 现在的位置
	keyPos := findPos	// 上一个key的位置
	maxLen := lenTarget
	byteArr := make([][]byte, 0)
	for i:=findPos; i<maxLen; i++ {
		if target[i] == '"' {
			isInString = !isInString
		} else if target[i] == '[' || target[i] == '{' {
			if !isInString {
				leftCount++
			}
		} else if target[i] == ']' || target[i] == '}' {
			if !isInString {
				leftCount--
			}
		} else if target[i] == ',' && leftCount == 1 {
			if !isInString {
				byteArr = append(byteArr, target[keyPos:i])
				keyPos=i+1
			}
		}
	}

	// 最后一个元素不可漏掉
	byteArr = append(byteArr, target[keyPos:maxLen-1])
	return byteArr
}


// 解析Map
func parseMap(target []byte) map[string] []byte{
	result := make(map[string][]byte, 0)
	lenOfTarget := len(target)
	if lenOfTarget <= 2 {
		return result
	}
	if target[0] != '{' || target[lenOfTarget-1] != '}' {
		return result
	}
	isInString := false	// 是否在引号内部
	findPos := 1		// 现在的位置
	leftCount := 1
	keyPos := findPos	// 上一个key的位置
	MaoPos := 0			// 上一个冒号的位置
	maxLen := len(target)
	for i:=findPos; i < maxLen; i++ {
		if target[i] == '{' || target[i] == '[' {
			if !isInString {
				leftCount++
			}
		} else if target[i] == '}' || target[i] == ']' {
			if !isInString {
				leftCount--
			}
		} else if target[i] == ':' && leftCount == 1 {
			if !isInString {	// 引号内的冒号不算
				MaoPos = i		// 设置冒号的位置
			}
		} else if target[i] == ',' && leftCount == 1 {
			if !isInString {	// 引号内的逗号不算
				result[string(target[keyPos+1:MaoPos-1])] = target[MaoPos+1:i]
				keyPos = i+1
			}
		} else if target[i] == '"' && (i == 0 || string(target[i-1]) != `\`) {
			isInString = !isInString
		}
	}
	// 最后一个元素
	result[string(target[keyPos+1:MaoPos-1])] = []byte(strings.Replace(string(target[MaoPos+1:maxLen-1]), `\"`, `"`, -1))
	return result
}


// 获取指定内容的范围
func getValueBetween(subByte []byte, finder []byte) (int, int) {

	var lenOfFinder = len(finder)
	var beginBytes = 0
	var endBytes = 0
	for i:=0; i+lenOfFinder<len(subByte); i++ {
		if bytes.Equal(subByte[i:i+lenOfFinder], finder) && subByte[i-1] == '"' && subByte[i+lenOfFinder] == '"' {
			if subByte[i+lenOfFinder+1] != ':' {
				return 0, 0
			}
			beginBytes = i-1
			switch subByte[i+lenOfFinder+2] {
			case '"':
				endBytes = getLastStringEnd(subByte[i+lenOfFinder+3:])
			}
		}
	}
	return beginBytes, beginBytes+endBytes+lenOfFinder+5
}


// 获取string的结束
func getLastStringEnd(target []byte) int {
	for i:=0; i<len(target); i++ {
		if target[i] == '"' && (i == 0 || string(target[i-1]) != `\`) {
			return i
		}
	}
	return 0
}


func findLastValue(subByte []byte, finder []byte) []byte {
	var lenOfFinder = len(finder)
	if lenOfFinder == 0 {
		return nil
	}
	for i:=0; i+lenOfFinder<len(subByte); i++ {
		if bytes.Equal(subByte[i:i+lenOfFinder], finder) && subByte[i-1] == '"' && subByte[i+lenOfFinder] == '"' {
			if subByte[i+lenOfFinder+1] != ':' {
				return nil
			}
			switch subByte[i+lenOfFinder+2] {
			case '"':
				return getLastString(subByte[i+lenOfFinder+3:])
			case '{':
				return getLastMap(subByte[i+lenOfFinder+2:])
			case '[':
				return getLastArr(subByte[i+lenOfFinder+2:])
			default:
				return getLastNumber(subByte[i+lenOfFinder+2:])
			}

			return []byte{}
		}
	}
	return nil
}


// 获取一个字符串
func getLastString(beginBytes []byte) []byte {
	for i:=0; i<len(beginBytes); i++ {
		if beginBytes[i] == '"' && (i == 0 || string(beginBytes[i-1]) != `\`) {
			return beginBytes[:i]
		}
	}
	return nil
}

// 获取一个数字
func getLastNumber(beginBytes []byte) []byte {
	for i:=0; i<len(beginBytes); i++ {
		if beginBytes[i] == ',' || beginBytes[i] == '}' || beginBytes[i] == ']' {
			return beginBytes[:i]
		}
	}
	//fmt.Printf("%v\n", string(beginBytes))
	return beginBytes
}


// 获取一个数组
func getLastArr(beginBytes []byte) []byte {
	leftCount := 1		// 左括号有几个
	for i:=1; i<len(beginBytes); i++ {
		if beginBytes[i] == '[' {
			leftCount++
		} else if beginBytes[i] == ']' {
			leftCount--
			if leftCount == 0 {
				return beginBytes[:i+1]
			}
		}
	}
	return nil
}

// 获取一个MAP
func getLastMap(beginBytes []byte) []byte {
	leftCount := 1
	for i:=1; i<len(beginBytes); i++ {

		if beginBytes[i] == '{' {
			leftCount++
		}else if beginBytes[i] == '}' {
			leftCount--
			if leftCount == 0 {
				return beginBytes[:i+1]
			}
		}
	}
	return nil
}



func JCode(code uint32, param string, v interface{}) []byte {

	if v == nil {
		return []byte(fmt.Sprintf("{\"msg\":%v,\"param\":\"%v\"}", code, param))
	}

	knd := reflect.ValueOf(v).Kind()

	switch knd {
	case reflect.String:
		return []byte(fmt.Sprintf("{\"msg\":%v,\"param\":\"%v\", \"data\":\"%v\"}", code, param, v.(string)))
	case reflect.Float64:
		fallthrough
	case reflect.Float32:
		return []byte(fmt.Sprintf("{\"msg\":%v,\"param\":\"%v\", \"data\":%0.2f}", code, param, v))
	case reflect.Int:
		fallthrough
	case reflect.Bool:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return []byte(fmt.Sprintf("{\"msg\":%v,\"param\":\"%v\", \"data\":%v}", code, param, v))
	default:
		jsonBytes := CreateBy(v)
		if jsonBytes == nil {
			return []byte(fmt.Sprintf("{\"msg\":%v,\"param\":\"%v\",\"data\":[]}", code, param))
		}
		return []byte(fmt.Sprintf("{\"msg\":%v,\"param\":\"%v\",\"data\":%v}", code, param, jsonBytes.ToStr()))
	}
}


func CsvDump(title []string, contents [] map[string] interface{}) []byte {

	resp := strings.Builder{}
	resp.WriteString(strings.Join(title, ","))
	for _, content := range contents {
		resp.WriteString("\n")
		lines := strings.Builder{}
		for _, val := range title {
			if lines.Len() > 0 {
				lines.WriteString(",")
			}
			getVal, exists := content[val]
			if !exists {
				lines.WriteString("null")
				continue
			}
			switch getVal.(type) {
			case float64:
				lines.WriteString(fmt.Sprintf("%0.2f", getVal))
			case float32:
				lines.WriteString(fmt.Sprintf("%0.2f", getVal))
			default:
				lines.WriteString(fmt.Sprintf("%v", getVal))
			}

		}
		resp.WriteString(lines.String())
	}
	return []byte(resp.String())
}