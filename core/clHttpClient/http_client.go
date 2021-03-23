package clHttpClient

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"clws-framework/core/clAes"
	"clws-framework/core/cljson"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)


// client
type ClHttpClient struct {
	url string					// 请求地址
	proxy string				// 设置代理
	method string				// 请求方式
	timeout uint32				// 设置超时时间
	encodetype uint32			// 加密方式
	aesKey string				// AES加密key
	simpleKey string			// 一般加密的key
	params map[string]string 	// 参数列表
	header map[string]string	// 请求头
	contentType uint32			// 请求文档类型
}

// 加密方式
const (
	ENCODE_TYPE_NONE = 0		// 明文
	ENCODE_TYPE_SIMPLE = 1		// 一般加密
	ENCODE_TYPE_AES = 2			// AES加密
)

const (
	ContentTypeForm = 0			// 正常form提交
	ContentParam = 1			// 参数提交(只允许GET)
	ContentJson = 2				// 通过json提交
)


// 获取一个新的对象
func NewClient(_url string) *ClHttpClient {

	client := ClHttpClient{
		url:    _url,
		method: "POST",
		timeout: 30,
		params: make(map[string]string),
		header: make(map[string]string),
	}

	client.header["Content-Type"] = "application/x-www-form-urlencoded"

	return &client
}


// 设置代理
func (this *ClHttpClient) SetProxy(_proxy string) {
	this.proxy = _proxy
}


// 设置超时时间
func (this *ClHttpClient) SetTimeout(_timeout uint32) {
	this.timeout = _timeout
}


// 设置方式
func (this *ClHttpClient) SetMethod(_method string) {
	this.method = _method
}


// 添加参数
func (this *ClHttpClient) AddParam(_key string, _val string) {

	this.params[_key] = _val
}


// 设置请求类型
func (this *ClHttpClient) SetContentType(_type uint32) {
	if _type == ContentJson {
		this.method = "POST"
	}
	this.contentType = _type
}


// 添加头
func (this *ClHttpClient) AddHeader(_key string, _val string) {

	this.header[_key] = _val
}


// 返回最终请求地址
func (this *ClHttpClient) Try() string {
	return this.BuildParamList()
}


//  开启AES加密
func (this *ClHttpClient) EnableAes(_aesKey string) {
	this.encodetype = ENCODE_TYPE_AES
	this.aesKey = _aesKey
}


// 开启一般加密
func (this *ClHttpClient) EnableSimpleEncode(_simpleKey string) {
	this.encodetype = ENCODE_TYPE_SIMPLE
	this.simpleKey = _simpleKey
}


// 构建参数
func (this *ClHttpClient) BuildParamList() string {

	// 参数拼接
	param_str := strings.Builder{}
	for PKey, PVal := range this.params {
		if param_str.Len() > 0 {
			param_str.WriteString("&")
		}
		param_str.WriteString(fmt.Sprintf("%v=%v", PKey, PVal))
	}

	var http_url = this.url

	if param_str.Len() == 0 {
		return http_url
	}

	return this.url + "?" + param_str.String()
}


// 开始请求
func (this *ClHttpClient) Do() (string, error) {

	var client *http.Client
	var proxyUrl *url.URL
	var err error

	if this.proxy != "" {

		proxyUrl, err = url.Parse(this.proxy)
		if err != nil {
			return "", err
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second * time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
				Proxy:           http.ProxyURL(proxyUrl),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
			},
		}
	}

	var http_url = ""
	var body io.Reader = nil
	if this.method == "POST" {
		http_url = this.url

		if this.contentType == ContentTypeForm {
			var r = http.Request{}
			r.ParseForm()
			bodyStr := ""
			for key, val := range this.params {
				r.Form.Add(key, val)
			}
			bodyStr = strings.TrimSpace(r.Form.Encode())
			body = strings.NewReader(bodyStr)
		} else if this.contentType == ContentJson {

			var jsonObj = cljson.M{}
			for key, val := range this.params {
				jsonObj[key] = val
			}
			body = strings.NewReader(cljson.CreateBy(jsonObj).ToStr())
		}

	}
	http_url = this.BuildParamList()

	req, err := http.NewRequest(this.method, http_url, body)
	if err != nil {

		return "", errors.New(fmt.Sprintf("HttpProxy: %v 请求地址: %v 错误:%v", this.proxy, http_url, err))
	}

	// 添加头
	for HKey, HVal := range this.header {
		req.Header.Add(HKey, HVal)
	}

	if this.contentType == ContentTypeForm {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if this.contentType == ContentJson {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("HttpProxy: %v 请求地址: %v 错误:%v", this.proxy, http_url, err))
	}

	jsonStr, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err2 != nil {
		return "", errors.New(fmt.Sprintf("HttpProxy: %v 请求地址: %v 错误:%v", this.proxy, http_url, err))
	}

	return string(jsonStr), nil
}


// 开始请求对象
func (this *ClHttpClient) DoJson(_iter interface{}) error {

	var client *http.Client
	var proxyUrl *url.URL
	var err error

	if this.proxy != "" {

		proxyUrl, err = url.Parse(this.proxy)
		if err != nil {
			return err
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second * time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
				Proxy:           http.ProxyURL(proxyUrl),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
			},
		}
	}


	// 参数拼接
	var http_url = this.BuildParamList()

	req, err := http.NewRequest(this.method, http_url, nil)
	if err != nil {
		return err
	}

	if this.contentType == ContentTypeForm {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if this.contentType == ContentJson {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// 添加头
	for HKey, HVal := range this.header {
		req.Header.Set(HKey, HVal)
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("URL: %v ERR: %v", http_url, err))
	}

	jsonStr, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err2 != nil {
		return errors.New(fmt.Sprintf("URL: %v ERR: %v", http_url, err2))
	}

	unmarshaErr := json.Unmarshal([]byte(jsonStr), _iter)
	if unmarshaErr != nil {
		return unmarshaErr
	}

	return nil
}


// 开始请求对象
func (this *ClHttpClient) DoXml(_xmldoc string) (error, string) {

	var client *http.Client
	var proxyUrl *url.URL
	var err error

	if this.proxy != "" {

		proxyUrl, err = url.Parse(this.proxy)
		if err != nil {
			return err, ""
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second * time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
				Proxy:           http.ProxyURL(proxyUrl),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
			},
		}
	}


	// 参数拼接
	var http_url = this.BuildParamList()

	req, err := http.NewRequest(this.method, http_url, strings.NewReader(_xmldoc))
	if err != nil {
		return err, ""
	}

	// 添加头
	for HKey, HVal := range this.header {
		req.Header.Set(HKey, HVal)
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("URL: %v ERR: %v", http_url, err)), ""
	}

	jsonStr, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err2 != nil {
		return errors.New(fmt.Sprintf("URL: %v ERR: %v", http_url, err2)), ""
	}


	return nil, string(jsonStr)
}



// 开始请求对象
func (this *ClHttpClient) DoJsonEx(_jsondoc string) (error, string) {

	var client *http.Client
	var proxyUrl *url.URL
	var err error

	if this.proxy != "" {

		proxyUrl, err = url.Parse(this.proxy)
		if err != nil {
			return err, ""
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second * time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
				Proxy:           http.ProxyURL(proxyUrl),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
			},
		}
	}


	// 参数拼接
	var http_url = this.BuildParamList()

	req, err := http.NewRequest(this.method, http_url, strings.NewReader(_jsondoc))
	if err != nil {
		return err, ""
	}

	// 添加头
	for HKey, HVal := range this.header {
		req.Header.Set(HKey, HVal)
	}

	if this.contentType == ContentTypeForm {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if this.contentType == ContentJson {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}


	res, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("URL: %v ERR: %v", http_url, err)), ""
	}

	jsonStr, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err2 != nil {
		return errors.New(fmt.Sprintf("URL: %v ERR: %v", http_url, err2)), ""
	}


	return nil, string(jsonStr)
}


// 开始
func (this *ClHttpClient) DoEncryptJson(_key string, _iter interface{}) error {

	var client *http.Client
	var proxyUrl *url.URL
	var err error

	if this.proxy != "" {

		proxyUrl, err = url.Parse(this.proxy)
		if err != nil {
			return err
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second * time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
				Proxy:           http.ProxyURL(proxyUrl),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
			},
		}
	}


	// 参数拼接
	param_str := strings.Builder{}
	for PKey, PVal := range this.params {
		if param_str.Len() > 0 {
			param_str.WriteString("&")
		}
		param_str.WriteString(fmt.Sprintf("%v=%v", PKey, PVal))
	}

	var http_url = this.url
	if param_str.Len() > 0 {
		http_url += "?p=" + string( clAes.Encode([]byte(param_str.String()), []byte(_key)) )
	}

	req, err := http.NewRequest(this.method, http_url, nil)
	if err != nil {
		return err
	}

	// 添加头
	for HKey, HVal := range this.header {
		req.Header.Set(HKey, HVal)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	jsonStr, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err2 != nil {
		return err
	}

	if len(jsonStr) == 0 {
		return nil
	}

	if string(jsonStr[0]) != "{" && string(jsonStr[0]) != "[" {
		jsonStr = clAes.Decode(jsonStr, []byte(_key))
	}

	unmarshaErr := json.Unmarshal(jsonStr, _iter)
	if unmarshaErr != nil {
		fmt.Printf(">> Unmarshal Err: %v \n (%v)\n", unmarshaErr, string(jsonStr))
		return err
	}

	return nil
}