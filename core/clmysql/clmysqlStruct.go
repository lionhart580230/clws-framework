package clmysql

type TdbResult map[string] string
type TmapResult map[string] TdbResult


/**
	DB数据集
 */
type DbResult struct {
	ArrResult [] TdbResult
	MapResult map[string] TdbResult
	Length uint32
}

func (res *DbResult) GetLength() uint32 {
	return res.Length
}


// 转变为map[string]interface{}结构
func (this *TdbResult) ToInterface() map[string] interface{} {

	resp := make(map[string] interface{})
	for key, val := range *this {
		resp[key] = val
	}
	return resp
}