package clredis

import (
	"agit.bxvip.co/golang-api/share-ciaolan/cljson"
	"testing"
	"fmt"
	"time"
)

var rd *RedisObject
func init() {
	rd, _ = New("127.0.0.1:6379","")
}

func TestRedisObject_Set(t *testing.T) {
	//err := rd.Set("TEST", "123456", 30)
	//fmt.Printf("写入失败:%v\n", err)

	err := rd.Set("TEST", cljson.M{
		"uid": 100,
		"username": "ok",
	}, 600)
	fmt.Printf("写入失败:%v\n", err)
}

func TestRedisObject_Get(t *testing.T) {
	fmt.Printf("|%v|", rd.GetJson("TEST").ToStr())
}

func TestRedisObject_HSet(t *testing.T) {
	rd.HSet("getSysInfo", "TESTING", cljson.M{"aa":000},20)
	<-time.After(time.Second*1)
	fmt.Printf(">>\n%v\n", rd.HGet("getSysInfo", "TESTING"))
}


func TestRedisObject_SetEx(t *testing.T) {

	for i:=0;i<10;i++ {
		go rd.SetEx("asdasd002", "qweqwe", 5)
	}
	
	<-time.After(10*time.Second)

}


func TestRedisObject_SetNx(t *testing.T) {

	for i:=0;i<100000;i++ {
		ok := rd.SetNx("testkey", "1", 10)
		fmt.Printf(">> %v\n", ok)
	}

}


func TestRedisObject_DelAll(t *testing.T) {
	rd.DelAll("*map*")
}


func TestRedisObject_IsExists(t *testing.T) {

	fmt.Printf("exists: %v\n", rd.IsExists("asd"))
}



func TestRedisObject_Increment(t *testing.T) {

	fmt.Printf("increment: %v\n", rd.Increment("a", -1000))

}


func TestRedisObject_Lpush(t *testing.T) {
	//rd.Lpush("axx", 86400, "2")

	fmt.Printf("POP: %v\n", rd.Rpop("axx"))
}


func TestRedisObject_GetKeys(t *testing.T) {

	fmt.Printf("KEYS: %+v\n", rd.GetKeys("*ONLINE*"))

}

func TestRedisObject_SetNXInt(t *testing.T) {

	fmt.Printf("OK: %v\n", rd.SetNXInt("I18NGUID", 10))
	fmt.Printf("CUR: %v\n", rd.Increment("I18NGUID", 1))
}