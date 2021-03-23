package clmysql

import (
	"clws-framework/core/clDebug"
	"testing"
	"time"
)


var DB *DBPointer
// 压测
func init() {
	mDB, err := NewDB("127.0.0.1", "root", "root", "shop")

	if err != nil {
		clDebug.Err("连接数据库失败! 错误: %v", err)
		return
	}

	DB = mDB
}

type GameTouzhuLog struct {
	Guid uint64 `db:"guid"`
	Username string `db:"username"`
	Price float64 `db:"price"`
}

func TestQuery(t *testing.T) {

	now := time.Now()
	res, err := DB.Query("select guid, username, price from game_touzhu_log order by guid desc limit 50000");
	if err != nil {
		clDebug.Err("查询投注记录失败! 错误: %v", err)
		return
	}

	clDebug.Debug("原始方式查询:%v条投注记录成功! 耗时: %0.2f秒", res.Length, time.Since(now).Seconds())


	now = time.Now()
	var GameTouzhuLog = make([]GameTouzhuLog, 0)

	err = DB.NewBuilder().Table("game_touzhu_log").Order("guid desc").Limit(0, 50000).FindAll(&GameTouzhuLog)
	if err != nil {
		clDebug.Err("查询投注记录失败! 错误: %v", err)
		return
	}

	clDebug.Debug("Object方式查询:%v条投注记录成功! 耗时:%0.2f秒", len(GameTouzhuLog), time.Since(now).Seconds())
}

func TestSqlBuider_Select(t *testing.T) {
	DB.Query("SELECT fish,fruit,poker,slot FROM request_time WHERE tag = 'fg'")
}


func TestSqlBuider_AddObj(t *testing.T) {

	DB.NewBuilder().Table("game_user").OnDuplicateKey([]string{
		"username",
		"Price",
	}).Add(map[string]interface{}{
		"guid":     1,
		"username": "asdasd001",
		"price":    100,
	})

}