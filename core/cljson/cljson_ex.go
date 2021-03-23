package cljson

import (
	"fmt"
	"reflect"
)

// 删除指定的key
func (this *JsonStream)Del() bool {
	begin, end := this.lastLockBegin, this.lastLockEnd
	if begin == 0 && end == 0 {
		return false
	}

	for i := begin-1; i >=0; i-- {
		if this.data[i] == '{' {
			if this.data[end] == ',' {
				end++
			}
			begin = i+1
			break
		} else if this.data[i] == '[' {
			if this.data[end] == ',' {
				end++
			}
			begin = i+1
			break
		} else if this.data[i] == ',' {
			begin = i
			break
		}
	}

	this.data = []byte(string(this.data[:begin])+ string(this.data[end:]))
	this.lastLockBegin, this.lastLockEnd = 0, 0
	return true
}


// 设置指定的key为指定的Str
func (this *JsonStream)SetStr(val string) {
	if this.lastLockBegin == 0 && this.lastLockEnd == 0{
		return
	}

	for i := this.lastLockBegin-1; i >=0; i-- {
		if this.data[i] == ':' {
			this.lastLockBegin = i+1
			break
		}
	}

	this.data = []byte(string(this.data[:this.lastLockBegin]) + `"`+val+`"` + string(this.data[this.lastLockEnd:]) )
	this.lastLockBegin, this.lastLockEnd = 0, 0
}

// 设置指定的key
func (this *JsonStream) SetObject(val interface{}) {
	if this.lastLockEnd == 0 {
		return
	}
	for i := this.lastLockBegin-1; i >=0; i-- {
		if this.data[i] == ':' {
			this.lastLockBegin = i+1
			break
		}
	}

	addStr := ""
	if val == nil {
		addStr = "null"
	} else {
		switch reflect.ValueOf(val).Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.Map:
			addStr = CreateBy(val).ToStr()
		case reflect.String:
			addStr = fmt.Sprintf("\"%v\"", val)
		default:
			addStr = fmt.Sprintf("%v", val)
		}
	}
	this.data = []byte(string(this.data[:this.lastLockBegin]) + addStr + string(this.data[this.lastLockEnd:]) )
	this.lastLockBegin, this.lastLockEnd = 0, 0
}

// 设置指定的key
func (this *JsonStream)SetValues(val interface{}) {
	if this.lastLockBegin == 0 && this.lastLockEnd == 0{
		return
	}

	for i := this.lastLockBegin-1; i >=0; i-- {
		if this.data[i] == ':' {
			this.lastLockBegin = i+1
			break
		}
	}

	this.data = []byte(string(this.data[:this.lastLockBegin]) + fmt.Sprintf("%v",val) + string(this.data[this.lastLockEnd:]) )
	this.lastLockBegin, this.lastLockEnd = 0, 0
}


// 递归查找某个key的区间
func (this *JsonStream) GetKey(key string, subkey ...string) (*JsonStream) {
	b, e := 0, 0
	if this.lastLockBegin == 0 && this.lastLockEnd == 0 {
		b,e = GetJsonValueEx(this.data, []byte(key))
	} else {
		b,e = GetJsonValueEx(this.data[this.lastLockBegin:this.lastLockEnd], []byte(key))
	}

	b, e = b+ this.lastLockBegin, e+this.lastLockBegin
	for _, val := range subkey {

		bt, et := GetJsonValueEx(this.data[b:e], []byte(val))
		if bt == 0 && et == 0 {
			return this
		}
		b, e = bt+b, et+b
	}

	this.lastLockBegin, this.lastLockEnd = b, e
	return this
}

// 重置
func (this *JsonStream) ResetOffset() {
	this.lastLockBegin = 0
	this.lastLockEnd = 0
}


// 递归查找某个key的区间
func (this *JsonStream) GetOffset(offset int) (*JsonStream){

	if this.lastLockBegin == 0 && this.lastLockEnd == 0 {
		this.lastLockBegin = 0
		this.lastLockEnd = len(this.data)
	}

	begin, end := this.lastLockBegin, this.lastLockEnd

	nowOffset := 0
	lastBegin := begin+1
	finallyEnd := 0

	leftKuo := 0
	leftFang := 0
	inString := false
	for i:=lastBegin;i<end;i++ {
		if this.data[i] == '[' {
			leftFang++
		} else if this.data[i] == '{' {
			leftKuo++
		} else if this.data[i] == '"' {
			inString = !inString
		} else if this.data[i] == ']' {
			leftFang--
		} else if this.data[i] == '}' {
			leftKuo--
		}

		if leftFang > 0 || leftKuo > 0 || inString {
			continue
		}

		if this.data[i] == ',' || this.data[i] == ']' {
			nowOffset ++

			if nowOffset > offset {
				finallyEnd = i
				break
			} else {
				lastBegin = i+1
			}
		}
	}

	this.lastLockBegin, this.lastLockEnd = begin+lastBegin, finallyEnd
	return this
}


// 获取JsonValue的起始和结尾
func GetJsonValueEx(target []byte, key []byte) (/* begin*/ int, /* end */int) {
	lenOfTarget := len(target)		// 想要搜寻的字节集
	lenOfKey := len(key)			// key的长度
	offset := 0						// 内容偏移

	objType := JSON_TYPE_STR

	lessObject := 0
	lessArray := 0
	inString := false

	switch target[0] {
	case 123:
		objType = JSON_TYPE_MAP
	case 91:
		objType = JSON_TYPE_ARR
	}

	if objType == JSON_TYPE_STR {
		return 0, 0
	}

	for i:=1; i<lenOfTarget-lenOfKey-3; i++ {

		switch target[i] {
		case 34:
			if target[i-1] != 92 {
				inString = !inString
			}
		case 123:
			if !inString {
				lessObject++
			}
		case 91:
			if !inString {
				lessArray++
			}
		case 125:
			if !inString {
				lessObject--
			}
		case 93:
			if !inString {
				lessArray--
			}
		}

		if lessObject != 0 || lessArray != 0 {
			continue
		}

		if inString && string(target[i:i+lenOfKey+3]) == `"`+string(key)+`":` {
			offset = i+lenOfKey+3
			break
		}
	}

	if offset == 0 {
		// 未找到
		return 0, 0
	}

	endKey := ','
	objectBet := 0
	ArrayBet := 0
	switch target[offset] {
	case 34:			// 字串
		endKey = '"'
	case '{':			// 对象
		endKey = '}'
		objectBet=1
	case '[':			// 数组
		endKey = ']'
		ArrayBet=1
	}

	for i:=offset+1;i<lenOfTarget;i++ {

		if endKey == '"' {
			if target[i] == uint8('"') {
				if target[i-1] != '\\' {
					return offset, i+1
				}
			}
			continue
		}

		if (target[i] == '}' || target[i] == ']') && endKey == ',' {
			return offset, i
		}

		if target[i] == uint8('}') {
			objectBet--
		} else if target[i] == uint8(']') {
			ArrayBet--
		} else if target[i] == uint8('{') {
			objectBet++
		} else if target[i] == uint8('[') {
			ArrayBet++
		}


		if target[i] == uint8(endKey) {
			if objectBet == 0 && ArrayBet == 0 {
				if endKey == 44 {
					return offset, i
				} else {
					return offset, i+1
				}
			}
		}
	}

	//fmt.Printf(">> objectBet %v , ArrayBet: %v\n", objectBet, ArrayBet)
	return 0, 0
}


// 删除指定的索引
// @param offset int 索引id,从0开始
func (this *JsonStream)DelOffset(offset int) bool {
	begin, end := this.lastLockBegin, this.lastLockEnd
	if begin == 0 && end == 0 {
		return false
	}

	if this.data[begin] != '[' || this.data[end-1] != ']' {
		return false
	}

	nowOffset := 0
	lastBegin := begin+1
	finallyEnd := 0
	for i:=lastBegin;i<end;i++ {
		if this.data[i] == ',' {
			nowOffset++

			if nowOffset > offset {
				finallyEnd = i+1
				break
			} else {
				lastBegin = i+1
			}
		}else if this.data[i] == ']' {
			nowOffset ++

			if nowOffset > offset {
				finallyEnd = i
				if this.data[lastBegin-1] == ',' {
					lastBegin--
				}
				break
			} else {
				lastBegin = i+1
			}
		}
	}

	if finallyEnd == 0 {
		return false
	}
	//fmt.Printf(">> %v\n", string(this.data[lastBegin:finallyEnd]))
	//this.data = []byte(string(this.data[:begin]))
	this.data = []byte(string(this.data[:lastBegin])+ string(this.data[finallyEnd:]))
	return true
}


// 设置指定的位置为某个对象
func (this *JsonStream)SetOffset(offset int, value interface{}) {
	begin, end := this.lastLockBegin, this.lastLockEnd
	if begin == 0 && end == 0 {
		return
	}

	if this.data[begin] != '[' || this.data[end-1] != ']' {
		return
	}

	nowOffset := 0
	lastBegin := begin+1
	finallyEnd := 0
	for i:=lastBegin;i<end;i++ {
		if this.data[i] == ',' || this.data[i] == ']' {
			nowOffset ++

			if nowOffset > offset {
				finallyEnd = i
				break
			} else {
				lastBegin = i+1
			}
		}
	}

	if finallyEnd == 0 {
		return
	}

	addStr := ""
	if value == nil {
		addStr = "null"
	} else {
		switch reflect.ValueOf(value).Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.Map:
			//fmt.Printf(">> %v\n", reflect.ValueOf(value).Kind())
			addStr = CreateBy(value).ToStr()
		case reflect.String:
			addStr = fmt.Sprintf("\"%v\"", value)
		default:
			addStr = fmt.Sprintf("%v", value)
		}
	}

	this.data = []byte(string(this.data[:lastBegin]) + addStr + string(this.data[finallyEnd:]))
	//fmt.Printf(">> %v\n%v\n", addStr, string(this.data))
	return
}


func (this *JsonStream)ArrayPush(value interface{}) {
	begin, end := this.lastLockBegin, this.lastLockEnd
	if begin == 0 && end == 0 {
		return
	}

	if this.data[begin] != '[' || this.data[end-1] != ']' {
		return
	}

	addStr:=""
	if value == nil {
		addStr = "null"
	} else {
		switch reflect.ValueOf(value).Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.Map:
			//fmt.Printf(">> %v\n", reflect.ValueOf(value).Kind())
			addStr = CreateBy(value).ToStr()
		case reflect.String:
			addStr = fmt.Sprintf("\"%v\"", value)
		default:
			addStr = fmt.Sprintf("%v", value)
		}
	}


	if begin < end-2 {
		addStr = ","+addStr
	}

	this.data = []byte(string(this.data[:end-1]) + addStr + string(this.data[end-1:]))
	return
}