package clGlobal

import (
	"github.com/xiaolan580230/clhttp-framework/core/clmysql"
	"github.com/xiaolan580230/clhttp-framework/core/skyconfig"
	"github.com/xiaolan580230/clhttp-framework/core/skylog"
	"github.com/xiaolan580230/clhttp-framework/core/skyredis"

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
var mRedis *skyredis.RedisObject
var mMysql *clmysql.DBPointer
var conf *skyconfig.Config

func Init(_filename string) {

	conf = skyconfig.New(_filename, 0)

	conf.GetStr("mongodb", "mgo_url", "", &SkyConf.MgoUrl)
	conf.GetStr("mongodb", "mgo_dbname", "", &SkyConf.MgoDBName)
	conf.GetStr("mongodb", "mgo_user", "", &SkyConf.MgoUser)
	conf.GetStr("mongodb", "mgo_pass", "", &SkyConf.MgoPass)

	conf.GetStr("mysql", "mysql_host", "", &SkyConf.MysqlHost)
	conf.GetStr("mysql", "mysql_name", "", &SkyConf.MysqlName)
	conf.GetStr("mysql", "mysql_user", "", &SkyConf.MysqlUser)
	conf.GetStr("mysql", "mysql_pass", "", &SkyConf.MysqlPass)

	conf.GetStr("redis", "redis_host", "", &SkyConf.RedisHost)
	conf.GetStr("redis", "redis_prefix", "", &SkyConf.RedisPrefix)
	conf.GetStr("redis", "redis_password", "", &SkyConf.RedisPass)

	conf.GetUint32("system", "log_type", skylog.LOG_TYPE_ALI, &SkyConf.LogType)
	conf.GetUint32("system", "log_level", skylog.LOG_LEVEL_ALL, &SkyConf.LogLevel)

	conf.GetStr("system", "version", "", &ServerVersion)
	conf.GetBool("system", "is_cluster", false, &SkyConf.IsCluster)
	conf.GetBool("system", "debug_router", false, &SkyConf.DebugRouter)
	skylog.New(ServerVersion)

	skylog.LogDebug("%+v", SkyConf)
	skylog.SetLevel(SkyConf.LogLevel)
	skylog.SetType(SkyConf.LogType)
}


// 获取redis连线
func GetRedis() *skyredis.RedisObject {
	if mRedis != nil && mRedis.Ping() {
		return mRedis
	}
	newRedis, err := skyredis.New(SkyConf.RedisHost, SkyConf.RedisPass, SkyConf.RedisPrefix)
	if err != nil {
		skylog.LogErr("连接redis [%v] [%v] 失败! %v", SkyConf.RedisHost, SkyConf.RedisPass, err)
		return nil
	}
	mRedis = newRedis
	return mRedis
}


// 获取mysql连线
func GetMysql() *clmysql.DBPointer {
	if mMysql != nil && mMysql.IsUsefull() {
		return mMysql
	}

	db, err := clmysql.NewDB(SkyConf.MysqlHost, SkyConf.MysqlUser, SkyConf.MysqlPass, SkyConf.MysqlName)
	if err != nil {
		return nil
	}
	mMysql = db
	return mMysql
}


// 获取mongodb连线
func GetMongo() {

}