package clCommon

import (
	"regexp"
)


// 判断是否是国内、外手机号码
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckPhone(number string) bool {

	reg1, err1 := regexp.Match(`^(\+86\s)?(13|14|15|16|17|18|19)[0-9]{9}$`, []byte(number))
	if err1 != nil {
		return false
	}
	reg2, err2 := regexp.Match(`^+[0-9]{1,6}\s[0-9]{5,12}$`, []byte(number))
	if err2 != nil {
		return false
	}

	if reg1 || reg2 {
		return true
	}

	return false
}


// 判断是否是邮箱
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckEmail(email string) bool {
	reg, err := regexp.Match(`^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$`, []byte(email))
	if err != nil {
		return false
	}
	return reg
}

// 判断帐号格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckUsername(username string) bool {
	reg, err := regexp.Match(`^[0-9a-zA-Z\_\.]{4,20}$`, []byte(username))
	if err != nil {
		return false
	}
	return reg
}

// 判断密码格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckPassword(password string) bool {
	reg, err := regexp.Match(`^[0-9a-zA-Z]{4,20}$`, []byte(password))
	if err != nil {
		return false
	}
	return reg
}

// 判断URL格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckURL(url string) bool {
	reg, err := regexp.Match(`^([hH][tT]{2}[pP]:\/\/|[hH][tT]{2}[pP][sS]:\/\/)([a-zA-Z0-9\-\_\p{Han}\%\&\?\#\@\:\.\/])+[0-9a-zA-Z\#]+$`, []byte(url))
	if err != nil {
		return false
	}
	return reg
}

// 判断IP格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckIP(ip string) bool {
	reg, err := regexp.Match(`^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$`, []byte(ip))
	if err != nil {
		return false
	}
	return reg
}

// 判断域名格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckDomain(domain string) bool {
	reg, err := regexp.Match(`^[a-zA-Z0-9\-\_\p{Han}]+(\.[a-zA-Z0-9\-\p{Han}]+)+$`, []byte(domain))
	if err != nil {
		return false
	}
	return reg
}

// 判断银行卡格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckBankCard(bankcard string) bool {
	reg, err := regexp.Match(`^[0-9]{9,20}$`, []byte(bankcard))
	if err != nil {
		return false
	}
	return reg
}

// 判断日期格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckDate(date string) bool {
	reg, err := regexp.Match(`^[0-9]{4}\-([0][0-9]|[1][0-2])\-([0-2][0-9]|[3][0-1])$`, []byte(date))
	if err != nil {
		return false
	}
	return reg
}

// 判断时间日期格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckDateTime(datetime string) bool {
	reg, err := regexp.Match(`^[0-9]{4}\-([0][0-9]|[1][0-2])\-([0-2][0-9]|[3][0-1])\s([0-1][0-9]|[2][0-3])\:[0-5][0-9]\:[0-5][0-9]$`, []byte(datetime))
	if err != nil {
		return false
	}
	return reg
}

// 判断时间格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckTime(datetime string) bool {
	reg, err := regexp.Match(`^([0-1][0-9]|[2][0-3])\:[0-5][0-9]\:[0-5][0-9]$`, []byte(datetime))
	if err != nil {
		return false
	}
	return reg
}


// 判断QQ账号是否符合要求
// @return 匹配返回TRUE, 不匹配返回FALSE
func CheckQQ(qq string) bool {
	reg, err := regexp.Match(`^[1-9][0-9]{5,11}$`, []byte(qq))
	if err != nil {
		return false
	}
	return reg
}

// 判断微信账号是否符合要求
// @return 匹配返回TRUE, 不匹配返回FALSE
func CheckWechat(wechat string) bool {
	reg, err := regexp.Match(`^[0-9a-zA-Z\_]{4,20}$`, []byte(wechat))
	if err != nil {
		return false
	}
	return reg
}

// 判断图片文件名是否满足要求
// @return 匹配返回TRUE, 不匹配返回FALSE
func CheckImageName(image string) bool {
	reg, err := regexp.Match(`^[0-9a-zA-Z]{1,48}\.(jpg|png)$`, []byte(image))
	if err != nil {
		return false
	}
	return reg
}

/*
	判断是否是数字  针对混合类型的字符串
	@param param string
	@return true/false  bool
*/
func CheckNumber(param string)(bool){
	reg,err := regexp.Match(`^[0-9]+$`,[]byte(param))
	if err != nil {
		return false
	}
	return reg
}

// 判断订单号格式
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckOrder(order string) bool {
	reg, err := regexp.Match(`^[\w\s-]{10,40}$`, []byte(order))
	if err != nil {
		return false
	}
	return reg
}

// 判断支付宝账号,可以是邮箱和手机号
// @return 匹配返回TRUE，不匹配返回FALSE
func CheckAlipay(alipay string) bool {
	reg, err := regexp.Match(`^(?:\w+\.?)*\w+@(?:\w+\.)+\w+|\d{9,11}$`, []byte(alipay))
	if err != nil {
		return false
	}
	return reg
}


// 验证字符串是否包含不安全的字符
func CheckIsSafe(val string) bool {
	reg, _ := regexp.Match(`([\;\"\'\\\<\>\/])|(script)|(acript)`, []byte(val))
	return reg == false
}

// 验证字符串可否转为整数
func CheckIsInt(val string) bool {
	reg, _ := regexp.Match(`^(\-)?[0-9]{1,24}$`, []byte(val))
	return reg
}

// 验证字符串是否可转为小数类型
func CheckIsFloat(val string) bool {
	reg, _ := regexp.Match(`^(\-)?[0-9]{1,10}(\.[0-9]{1,4})?$`, []byte(val))
	return reg
}

// 验证是否验证码格式
func CheckIsVcode(val string) bool {
	reg, _ := regexp.Match(`^[0-9]{4}$`, []byte(val))
	return reg
}

// 验证是否客户端类型
func CheckIsClient(val string) bool {
	reg, _ := regexp.Match(`^[0-3]$`, []byte(val))
	return reg
}

// 验证是否是UUID
func CheckIsUUID(val string) bool {
	reg, _ := regexp.Match(`^[a-z0-9]{8}(\-[a-z0-9]{4}){3}\-[a-z0-9]{12}$`, []byte(val))
	return reg
}

// 验证是否是版本号格式
func CheckIsVersion(val string) bool {
	reg, _ := regexp.Match(`^v?[0-9]{1,3}(\.[0-9]{1,3}){1,3}$`, []byte(val))
	return reg
}

// 验证是否是一个图片base64格式
func CheckIsImageBase64(val string) bool {
	reg, _ := regexp.Match(`^data\-(jpg|png|gif)\-[a-zA-Z0-9\/\=\;\,\+]+$`, []byte(val))
	return reg
}