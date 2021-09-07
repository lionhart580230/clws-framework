package clGlobal

import (
	"github.com/xiaolan580230/clUtil/clConfig"
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clUtil/clMysql"
	"github.com/xiaolan580230/clUtil/clRedis"
)

var ServerVersion = `v1.0.0`

type SkyConfig struct {
	MgoUrl       string
	MgoDBName    string
	MgoUser      string
	MgoPass      string

	MysqlHost string
	MysqlName string
	MysqlUser string
	MysqlPass string

	RedisHost    string
	RedisPrefix  string
	RedisPass    string

	LogType  uint32
	LogLevel uint32

	IsCluster bool
	DebugRouter bool
}

var SkyConf SkyConfig
var mRedis *clRedis.RedisObject
var mMysql *clMysql.DBPointer

func Init(_filename string) {


	SkyConf.MgoUrl = clConfig.GetStr("mgo_url", "")
	SkyConf.MgoDBName = clConfig.GetStr("mgo_dbname", "")
	SkyConf.MgoUser = clConfig.GetStr("mgo_user", "")
	SkyConf.MgoPass = clConfig.GetStr("mgo_pass", "")

	SkyConf.MysqlHost = clConfig.GetStr("mysql_host", "")
	SkyConf.MysqlName = clConfig.GetStr("mysql_name", "")
	SkyConf.MysqlUser = clConfig.GetStr("mysql_user", "")
	SkyConf.MysqlPass = clConfig.GetStr("mysql_pass", "")

	SkyConf.RedisHost = clConfig.GetStr("redis_host", "")
	SkyConf.RedisPrefix = clConfig.GetStr("redis_prefix", "")
	SkyConf.RedisPass = clConfig.GetStr("redis_password", "")

	SkyConf.IsCluster = clConfig.GetBool("is_cluster", false)
	SkyConf.DebugRouter = clConfig.GetBool("debug_router", false)

	clLog.Debug("%+v", SkyConf)
}


// 获取redis连线
func GetRedis() *clRedis.RedisObject {
	if mRedis != nil && mRedis.Ping() {
		return mRedis
	}
	newRedis, err := clRedis.New(SkyConf.RedisHost, SkyConf.RedisPass, SkyConf.RedisPrefix)
	if err != nil {
		clLog.Error("连接redis [%v] [%v] 失败! %v", SkyConf.RedisHost, SkyConf.RedisPass, err)
		return nil
	}
	mRedis = newRedis
	return mRedis
}


// 获取mysql连线
func GetMysql() *clMysql.DBPointer {
	if mMysql != nil && mMysql.IsUsefull() {
		return mMysql
	}

	db, err := clMysql.NewDB(SkyConf.MysqlHost, SkyConf.MysqlUser, SkyConf.MysqlPass, SkyConf.MysqlName)
	if err != nil {
		return nil
	}
	mMysql = db
	return mMysql
}


// 获取mongodb连线
func GetMongo() {

}