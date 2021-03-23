package clHttpClient

import (
	"agit.bxvip.co/golang-api/share-ciaolan/common"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestClHttpClient_Do(t *testing.T) {


	httpClient := NewClient("https://zjgguolong.pandanb.com/Sapi")

	httpClient.AddParam("ac", "getBossStatic")
	httpClient.AddParam("key", "fuckfuck")
	httpClient.AddParam("bdate", "20190901")
	httpClient.AddParam("edate", "20190901")

	resp, err := httpClient.Do()
	fmt.Printf("Resp: %v\nCode: %v\n", resp, err)
}


type FixObj struct {
	Tag string `json:"tag"`
	Qishu uint64 `json:"qishu"`
	Ball string `json:"ball"`
}

func TestClHttpClient_DoJson(t *testing.T) {

	//users := []string{}
	for i:=0;i<=10;i++{
		//json := `"ac":"regUser","enom":"qawap.sg04.com","client_type":"1","username":"xiaoma0001","vid":"kzXFmFg0jJfKDG3SBRsT","qq":"","wechat":"","edition":"v1.0.0","password":"asdasd","vcode":"4103","email":"","tg_code":"p21","phone":"","real_name":""`
		//js := cljson.New([]byte(json))
		//data := js.ToMap().ToTree()
		//data["username"] = fmt.Sprintf("ciaolan00%v",i)
		//data["tg_code"] = fmt.Sprintf("p%v",common.RandInt(21,30))
		//users = append(users,  fmt.Sprintf("wanduzi00%v",i))
		//data["vid"] = "9F8DDDF6-5684-43FA-806D-05F1E850CF04"
		//data["vcode"] = "6666"
		////p := fmt.Sprintf("p=%v", AexEncode(data))
		////p = UrlEncode("?"+p)

		httpClient := NewClient("http://qawap.sg04.com/request")
		httpClient.AddParam("ac", "regUser")
		httpClient.AddParam("client_type", "1")
		httpClient.AddParam("enom", "qawap.sg04.com")
		httpClient.AddParam("username", fmt.Sprintf("ciaolan00%v",i))
		httpClient.AddParam("tg_code", fmt.Sprintf("p%v",common.RandInt(21,30)))
		httpClient.AddParam("password", "asdasd")
		httpClient.AddParam("edition", "v1.0.0")
		httpClient.AddParam("vid", "9F8DDDF6-5684-43FA-806D-05F1E850CF04")
		httpClient.AddParam("vcode", "6666")


		res, err := httpClient.Do()
		fmt.Println(res, err)
	}

}


func TestClHttpClient_AddParam(t *testing.T) {


	httpClient := NewClient("https://zjgguolong.pandanb.com/Sapi")
	httpClient.AddParam("ac", "getBossLotteryStatic")
	httpClient.AddParam("key", "1FB21322B989AC198B94D25C397BAEF9")
	httpClient.AddParam("bdate", "20200408")
	httpClient.AddParam("edate", "20200408")
	httpClient.AddParam("Timestamp", fmt.Sprintf("%v", uint32(time.Now().Unix())))
	httpClient.EnableAes("zLxJQczk1#1qS2Dz")

	resp, err := httpClient.Do()
	if err != nil {
		fmt.Printf( "获取业主[%v]报表数据失败! 错误: %v", "111",  err)
		return
	}
	fmt.Printf("param: %v\n", resp)

	//httpClient := NewClient("https://zjgguolong.pandanb.com/Sapi")
	//httpClient.AddParam("ac", "sadminGetAdminSysInfo")
	//httpClient.AddParam("key", "1FB21322B989AC198B94D25C397BAEF9")
	//httpClient.AddParam("Timestamp", fmt.Sprintf("%v", uint32(time.Now().Unix())))
	//
	//httpClient.EnableAes("zLxJQczk1#1qS2Dz")
	//fmt.Printf("param: %v\n", httpClient.Try())
	//resp, err := httpClient.Do()
	//fmt.Printf("resp: %v\n err:%v\n", resp, err)
}


func TestClHttpClient_SetTimeout(t *testing.T) {

	var r http.Request

	r.ParseForm()
	r.Form.Add("date", "2020")
	r.Form.Add("type", "1")

	bodystr := strings.TrimSpace(r.Form.Encode())
	request, _ := http.NewRequest("POST", "http://www.1988660.com/ci/Api/tw6/findSpeedSixHistory", strings.NewReader(bodystr))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf(">>%+v\n", *request)
	resp, _ := http.DefaultClient.Do(request)
	buffer, _ := ioutil.ReadAll( resp.Body )
	request.Body.Close()

	fmt.Printf(">>RESP: %v\n", string(buffer))
}