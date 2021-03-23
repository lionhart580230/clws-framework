package clmysql

import (
	"reflect"
	"strings"
	"fmt"
)


// 根据数据结构取得字段列表
func GetAllField(_val interface{}) []string {

	_fields := make([]string, 0)
	_type := reflect.TypeOf(_val).Elem()

	for i := 0; i < _type.NumField(); i++ {
		fieldName := _type.Field(i).Tag.Get("db")
		if fieldName == "" || fieldName == "-" {
			continue
		}
		_fields = append(_fields, fieldName)
	}
	return _fields
}


// 将数据反序列化为一个对象
func Unmarsha(_row TdbResult, _inter interface{}) {

	_type := reflect.TypeOf(_inter)
	_value := reflect.ValueOf(_inter)
	_valueE := _value.Elem()
	for i := 0; i < _value.Elem().NumField(); i++ {

		field_name := _type.Elem().Field(i).Tag.Get("db")

		switch _type.Elem().Field(i).Type.String() {
		case "uint32":
			_valueE.Field(i).SetUint(uint64(_row.GetUint32(field_name, 0)))
		case "uint64":
			_valueE.Field(i).SetUint(_row.GetUint64(field_name, 0))
		case "int32":
			_valueE.Field(i).SetInt(int64(_row.GetInt32(field_name, 0)))
		case "int64":
			_valueE.Field(i).SetInt(_row.GetInt64(field_name, 0))
		case "string":
			_valueE.Field(i).SetString(_row.GetStr(field_name, ""))
		case "bool":
			_valueE.Field(i).SetBool(_row.GetBool(field_name, false))
		case "float32":
			_valueE.Field(i).SetFloat(float64(_row.GetFloat32(field_name, 0)))
		case "float64":
			_valueE.Field(i).SetFloat(_row.GetFloat64(field_name, 0))
		}
	}
}


// 根据数据结构取得字段列表
func GetInsertSql(_val interface{}, _primary bool) ([]string, []string) {

	_fields := make([]string, 0)
	_values := make([]string, 0)

	_type := reflect.TypeOf(_val)
	_value := reflect.ValueOf(_val)

	for i := 0; i < _type.NumField(); i++ {

		if (!_primary) {
			if strings.ToUpper(_type.Field(i).Tag.Get("primary")) == "TRUE" {
				continue
			}
		}

		_fields = append(_fields, _type.Field(i).Tag.Get("db"))
		_values = append(_values, fmt.Sprintf("%v", _value.Field(i).Interface()))
	}
	return _fields, _values
}
