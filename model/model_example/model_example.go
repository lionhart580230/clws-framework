package model_example

import (
	"github.com/xiaolan580230/clws-framework/core/clDebug"
	"github.com/xiaolan580230/clws-framework/core/clGlobal"
)

const TableName = "tb_example"

type ModelExample struct {
	Id uint32 `db:"id",primary:"TRUE"`
	Name string `db:"name"`
}


// 获取数据
func GetData() []ModelExample {
	DB := clGlobal.GetMysql()
	if DB == nil {
		clDebug.Err("连接数据库失败!")
		return nil
	}

	data := make([]ModelExample, 0)

	err := DB.NewBuilder().Table(TableName).FindAll(&data)
	if err != nil {
		clDebug.Err("获取数据失败! 错误:%v", err)
		return nil
	}
	return data
}