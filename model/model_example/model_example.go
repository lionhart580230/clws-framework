package model_example

import (
	"github.com/xiaolan580230/clUtil/clLog"
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
		clLog.Error("连接数据库失败!")
		return nil
	}

	data := make([]ModelExample, 0)

	err := DB.NewBuilder().Table(TableName).FindAll(&data)
	if err != nil {
		clLog.Error("获取数据失败! 错误:%v", err)
		return nil
	}
	return data
}