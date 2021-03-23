### 新版数据库优化


> 使用对象方式读取数据库资料

~~~

// 定义一个数据库对象
// db字段为在数据库中的字段名称
type GameTouzhuLog struct {
	Guid uint64 `db:"guid"`
	Username string `db:"username"`
	Price float64 `db:"price"`
}

// 通过NewBuilder()创建一个数据库构建函数, 然后将数据放入 FindAll中，即可提取出整个数据库内容.
func TestQueryTx(t *testing.T) {

	now := time.Now()
	var GameTouzhuLog = make([]GameTouzhuLog, 0)

	err := DB.NewBuilder().Table("game_touzhu_log").Order("guid desc").Limit(0, 10000).FindAll(&GameTouzhuLog)
	if err != nil {
		cllog.Log.LogErr("", "查询投注记录失败! 错误: %v", err)
		return
	}

	cllog.Log.LogDebug("", "查询:%v条投注记录成功! 耗时:%0.2f秒", len(GameTouzhuLog), time.Since(now).Seconds())
}

~~~

除此以外还有: `FindOne` 函数进行获取单行数据, 如果数据不存在, 将会返回 `not found` 的错误.

~~~
// 引用上面的例子
// 由于是获取单行数据，所以不需要make出slice
var GameTouzhuInfo GameTouzhuLog
err := DB.NewBuilder().Table("game_touzhu_log").Order("guid desc").Limit(0, 10000).FindOne(&GameTouzhuInfo)
if err.Error() == "not found" {
   cllog.Log.LogErr("", "数据不存在!!")
   return
}

~~~