package clredis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/xiaolan580230/clws-framework/core/clCommon"
	"github.com/xiaolan580230/clws-framework/core/cljson"
	"strings"
	"sync"
	"time"
)

type RedisObject struct {
	myredis *redis.Client
	prefix string
	isCluster bool
}

var RedisPool map[string] *RedisObject
var Locker sync.RWMutex
func init() {
	RedisPool = make(map[string] *RedisObject)
}

func empty(data ...string )bool{
	if len(data) == 0 {
		return true
	}
	for _,v := range data {
		if v != "" {
			return false
		}
	}
	return true
}



func NewSimple(addr, website string) (*RedisObject, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		PoolSize: 10,
		PoolTimeout: 30 * time.Second,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	clrd := &RedisObject{
		myredis: client,
		prefix: website,
		isCluster:false,
	}

	return clrd, nil
}

func New(addr, web_site string) (*RedisObject, error) {

	Locker.RLock()
	val, find := RedisPool[addr+web_site]
	Locker.RUnlock()

	if find {
		redisPing := val.myredis.Ping()
		if redisPing.Err() == nil {
			return val, nil
		}
		Locker.Lock()
		delete(RedisPool, addr+web_site)
		Locker.Unlock()
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		PoolSize: 10,
		PoolTimeout: 30 * time.Second,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	clrd := &RedisObject{
		myredis: client,
		prefix: web_site,
		isCluster:false,
	}

	Locker.Lock()
	RedisPool[addr+web_site] = clrd
	Locker.Unlock()

	return clrd, nil
}


func (this *RedisObject)Close() {

	if this.myredis != nil {
		this.myredis.Close()
	}
}

// 删除
func (this *RedisObject) Del(key string) (error) {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}

	i := this.myredis.Del(keys)
	return i.Err()
}

// 设置
func (this *RedisObject) Set(key string, val interface{}, expire int32) error {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	err := this.myredis.Set(keys, buildRedisValue(keys, uint32(expire), val),
		time.Duration(time.Second * time.Duration(expire))).Err()
	return err
}


// 设置json
func (this *RedisObject)SetJson(key string, val interface{}, expire int32) error {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	err := this.myredis.Set(keys, cljson.CreateBy(buildRedisValue(keys, uint32(expire), val)).ToStr(),
		time.Duration(time.Second * time.Duration(expire))).Err()
	return err
}

// 获取指定的值
func (this *RedisObject)Get(key string) string {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	resp := this.myredis.Get(keys)
	result := checkRedisValid(keys, resp)
	if result == "" {
		this.myredis.Del(keys)
	}
	return result
}

func (this *RedisObject) GetNoPrefix(key string) string {

	keys := key
	resp := this.myredis.Get(keys)
	result := checkRedisValid(keys, resp)
	if result == "" {
		this.myredis.Del(keys)
	}
	return result
}


// 获取指定的json结构
func (this *RedisObject)GetJson(key string) *cljson.JsonStream {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	obj := this.myredis.Get(keys)
	return cljson.New([]byte(checkRedisValid(keys, obj)))
}

// 设置hash结构
func (this *RedisObject)HSet(key string, field string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	value = buildRedisValue(keys+field, expire, value)
	rest := this.myredis.HSet(keys, field, value)
	if rest == nil {
		return false
	}

	if _, err := rest.Result(); err != nil {
		fmt.Printf(">> HSet |%v|->|%v| Failed! Err:%v\n", keys, field, err)
		return false
	}

	return true
}


// 设置hash结构
func (this *RedisObject)SetEx(key string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	value = buildRedisValue(keys, expire, value)
	rest := this.myredis.HSet(keys, "test", value)
	if rest == nil {
		return false
	}

	ok, err := rest.Result()
	if !ok {
		fmt.Printf("写入失败: %v\n", ok)
		return false
	}

	if err != nil {
		fmt.Printf("写入失败: 错误:%v\n", err)
	}

	fmt.Printf(">> 写入成功!\n")
	return true
}




// 设置hash结构的值(保存为json)
func (this *RedisObject)HSetJson(key string, field string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	value = buildRedisValue(keys+field, expire, value)
	rest := this.myredis.HSet(keys, field, value)

	if rest == nil {
		return false
	}

	return rest.Val()
}

// 获取hash结构
func (this *RedisObject)HGet(key string, field string) string {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	resp := this.myredis.HGet(keys, field)
	result := checkRedisValid(keys + field, resp)
	if result == "" {
		this.myredis.HDel(keys, field)
	}
	return result
}

func (this *RedisObject)HDel(key string, field string) bool {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	resp := this.myredis.HDel(keys, field)
	return resp.Val() > 0
}

// 获取hash结构的值
func (this *RedisObject)HGetJson(key string, field string) *cljson.JsonStream {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	val := this.myredis.HGet(keys, field)
	res := checkRedisValid(keys + field, val)
	if res == "" {
		return nil
	}
	return cljson.New([]byte(res))
}

// 获取全部的key
func (this *RedisObject)HGetKeys(key string, prefix string) []string {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}

	val := this.myredis.HKeys(keys)
	if val == nil {
		return []string{}
	}
	resp := make([]string, 0)
	for _, val := range val.Val() {
		if strings.HasPrefix(val, prefix) {
			resp = append(resp, val)
		}
	}
	return resp
}

// 删除指定开头的keys
func (this *RedisObject)HDelKeys(key string, prefix string)  {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	keylist := this.HGetKeys(keys, prefix)
	if len(keylist) > 0 {
		this.myredis.HDel(keys, keylist...)
	}
}

// 获取全部的hash字段
func (this *RedisObject)HGetAll(key string) map[string] string {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}

	val := this.myredis.HGetAll(keys)
	return checkRedisValidMap(keys, val)
}


// 设置锁
func (this *RedisObject) SetNx(key string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	value = buildRedisValue(keys, expire, value)
	rest := this.myredis.SetNX(keys, value, time.Duration(expire)*time.Second)
	if rest == nil {
		return false
	}

	if _, err := rest.Result(); err != nil {
		fmt.Printf(">> SetNX |%v| Failed! Err:%v\n", keys, err)
		return false
	}

	return rest.Val()
}


// 检验redis缓存是否有效
// @param keys string redis缓存的键名
// @param targetData *StringCmd 目标数据
func checkRedisValidMap(keys string, targetData *redis.StringStringMapCmd) map[string] string {
	if targetData == nil || len(targetData.Val()) == 0 {
		return nil
	}

	resp := make(map[string] string)

	for key, val := range targetData.Val() {
		js := cljson.New([]byte(val))
		expireTime := js.GetUint32("expire")
		// 缓存到期
		if expireTime > 0 && expireTime < uint32(time.Now().Unix()) {
			continue
		}

		sign := clCommon.Md5("Cache:__"+keys+key)
		if js.GetStr("sign") != sign {
			continue
		}

		resp[key] = js.GetStr("data")
	}
	return resp
}


// 检验redis缓存是否有效
// @param keys string redis缓存的键名
// @param targetData *StringCmd 目标数据
func checkRedisValid(keys string, targetData *redis.StringCmd) string {
	if targetData == nil || targetData.Val() == "" {
		return ""
	}

	js := cljson.New([]byte(targetData.Val()))
	if js == nil {
		return ""
	}

	expireTime := js.GetUint32("expire")
	addtime := js.GetUint32("addtime")

	// 缓存到期  新增添加时间大于当前时间表示有问题
	if expireTime == 0 || expireTime < uint32(time.Now().Unix()) {
		return ""
	}
	if addtime > uint32(time.Now().Unix()) {
		return ""
	}

	sign := clCommon.Md5("Cache:__"+keys)
	if js.GetStr("sign") != sign {
		return ""
	}

	return js.GetStr("data")
}

// 组装缓存的值
func buildRedisValue(keys string, expire uint32, data interface{}) string {
	return cljson.CreateBy(cljson.M{
		"data": data,
		"addtime": uint32(time.Now().Unix()), // 写入redis缓存的时间
		"ip": "",
		"expire": uint32(time.Now().Unix()) + expire,
		"sign": clCommon.Md5("Cache:__"+keys),
	}).ToStr()
}


// 删除指定用户缓存
func (this *RedisObject) DelUserCache(uid uint32) {

	this.Del("USER"+"_"+string(uid))
}

// 删除指定接口缓存
// @param apiname string 接口名称
// @param uid uint32 用户的uid，如果为0则删除全部
func (this *RedisObject) DelApiCache(apiname string, uid uint32) {

	if uid == 0 {
		this.Del(apiname)
		return
	}
	this.HDelKeys(apiname, fmt.Sprintf("U%v_", uid))
}

// 操作list结构 lpush,push 是会更新key的过期时间
func (this *RedisObject)Lpush(key string,expire uint32,values ...interface{}) bool {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}
	for k,value:= range values {
		values[k] = buildRedisValue(keys, uint32(expire), value)
	}
	rest := this.myredis.LPush(keys, values...)

	// 设置过期时间
	if expire > 0 {
		this.myredis.Expire(keys, time.Duration(expire) * 1000 * time.Millisecond)
	}

	if rest == nil {
		return false
	}
	return true
}

// 操作list结构 lpop
func (this *RedisObject)Lpop(key string) string {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}

	val := this.myredis.LPop(keys)
	result := checkRedisValid(keys, val)
	return result
}


//取队列元素个数
func (this *RedisObject) Llen(key string) (error, int64) {
	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}

	result := this.myredis.LLen(keys)
	return result.Err(), result.Val()
}

// 操作list结构 rpop
func (this *RedisObject)Rpop(key string) interface{} {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}

	val := this.myredis.RPop(keys)
	result := checkRedisValid(keys, val)
	return result
}

// 删除list
func (this *RedisObject)DelList(key string) {

	keys := key
	if this.prefix != "" {
		keys = this.prefix+"_"+key
	}

	this.myredis.LTrim(keys, 1, 0)

	return
}


// 获取key列表
func (this *RedisObject)GetKeys(key string) []string {

	res := this.myredis.Keys(key)

	return res.Val()
}



// 删除所有的类似的key
func (this *RedisObject)DelAll(key string) {

	res := this.myredis.Keys(key)

	klist, _ := res.Result()
	this.myredis.Del(klist...)
}


// 判断键是否存在
func (this *RedisObject) IsExists(key string) bool {

	res := this.myredis.Exists(key)

	return res.Val() == 1
}


// 添加一个值
func (this *RedisObject) SetNXInt(key string, _val int64) bool {
	var res *redis.BoolCmd
	res = this.myredis.SetNX(key, _val, 0)
	return res.Val()
}


// 添加一个值
func (this *RedisObject) Increment(key string, _val int64) int64 {
	var res *redis.IntCmd
	if _val < 0 {
		res = this.myredis.DecrBy(key, -_val)
	} else {
		res = this.myredis.IncrBy(key, _val)
	}
	return res.Val()
}