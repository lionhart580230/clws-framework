package clGlobal

import (
	"clws-framework/core/clDebug"
	"clws-framework/core/clmysql"
	"clws-framework/core/clredis"
	"sync"
)

var mRedis *clredis.RedisObject
var mRedisLocker sync.RWMutex
var mMysql *clmysql.DBPointer
var mMysqlLocker sync.RWMutex

var mDBHost, mDBUser, mDBPass, mDBName string
var mRedisHost, mRedisPrefix string

// 初始化数据库配置
func InitMysqlConfig(_host, _user, _pass, _name string) {
	mDBHost = _host
	mDBUser = _user
	mDBPass = _pass
	mDBName = _name
}

// 初始化redis配置
func InitRedisConfig(_host, _prefix string) {
	mRedisHost = _host
	mRedisPrefix = _prefix
}


// 获取redis
func GetRedis() *clredis.RedisObject {
	mRedisLocker.Lock()
	defer mRedisLocker.Unlock()

	if mRedis == nil {
		if mRedisHost == "" {
			return nil
		}
		var err error
		mRedis, err = clredis.New(mRedisHost, mRedisPrefix)
		if err != nil {
			clDebug.Err("连接redis失败! 错误:%v", err)
			return nil
		}
 	}
 	return mRedis
}

// 获取数据库
func GetMysql() *clmysql.DBPointer {
	mMysqlLocker.Lock()
	defer mMysqlLocker.Unlock()

	if mMysql == nil {
		if mDBHost == "" {
			return nil
		}
		var err error
		mMysql, err = clmysql.NewDB(mDBHost, mDBUser, mDBPass, mDBName)
		if err != nil {
			clDebug.Err("连接mysql失败! 错误:%v", err)
			return nil
		}
	}
	return mMysql
}