package cljson


// 遍历
func (js *JsonArray)Each(hf func(key int, value *JsonStream) bool) {
	for k, val := range *js {
		if ok := hf(k, val.data); !ok {
			return
		}
	}
}

// 获取指定下标
func (js *JsonArray) GetOffset(offset int) (*JsonStream) {
	if offset >= len(*js) {
		return nil
	}
	return (*js)[offset].data
}

// 获取长度
func (js *JsonArray) GetLength() int {
	return len(*js)
}

// 获取标准的数组
func (js *JsonArray) ToCustom() []string {
	temp := make([]string, 0)
	js.Each(func (key int, val *JsonStream) bool {
		if val == nil {
			return true
		}
		temp = append(temp, val.ToStr())
		return true
	})
	return temp
}

func (js * JsonArray) IsEmpty() bool {
	if len(*js) == 0 {
		return true
	}
	return false
}