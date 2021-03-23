package cljson

const (
	JSON_TYPE_NIL = -1
	JSON_TYPE_ARR = 0
	JSON_TYPE_MAP = 1
	JSON_TYPE_STR = 2
	JSON_TYPE_BOOL = 3
	JSON_TYPE_INT = 4
	JSON_TYPE_NULL = 5

)

// 定义结构体
type JsonStream struct {
	data []byte    		// 数据 []byte
	dataLength uint32   // 数据长度
	dataType uint8		// 数据类型
	lastLockBegin int	// 这是什么
	lastLockEnd int		// 这是什么
	dataMap map[string] string
	dataArray []string
}


type JsonArray []jsonItem
type JsonMap map[string] jsonItem

type jsonItem struct {
	data *JsonStream
}


type JCodeStruct struct {
	Msg uint32
	Param string
	Data interface{}
}