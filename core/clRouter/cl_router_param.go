package clRouter

import (
	"github.com/xiaolan580230/clws-framework/core/clCommon"
	"net/url"
)



func NewParam(_key, _def string, _ptype int, _static bool) RouterParam {
	return RouterParam{
		Key:    _key,
		Def:    _def,
		PType:  _ptype,
		Static: _static,
	}
}


const (
	ParamTypeSafe = 0			// 只要不包含非法字符即可,不能为空
	ParamTypeInt = 1			// 必须是整数型
	ParamTypeUrl = 2		// 字符串，使用urlEncode进行转义
	ParamTypeHtml = 3		// 字符串，类似PHP中的htmlspecialchars进行转译
	ParamTypeFloat = 4		// 可以是整数，也可以带小数点
	ParamTypeAll = 5			// 不进行任何处理，直接放行(不推荐)
	ParamTypePhone = 6		// 手机号码格式
	ParamTypeEmail = 7		// 邮箱格式
	ParamTypeDomain = 8		// 域名格式
	ParamTypeIP = 9			// IP格式
	ParamTypeBankCard = 10    // 银行卡格式
	ParamTypeDate = 11		// 日期格式如: 2018-02-01
	ParamTypeDateTime = 12	// 时间日期格式如: 2018-02-01 00:00:00
	ParamTypeUserPass = 13	// 帐号或者密码的格式
	ParamTypeQQ = 14			// QQ账号
	ParamTypeWechat = 15		// 微信账号
	ParamTypeImageName = 16	// 图片文件名
	ParamTypeClientType = 17	// 设备ID
	ParamTypeVcode = 18		// 验证码格式
	ParamTypeUUID = 19		// UUID格式 如: ff1793b2-2825-11e8-b394-0242ac11000a
	ParamTypeVersion = 20		// 版本号。支持最多四级子版本
	ParamTypeImageBase64 = 24 // 图片的base64格式
	ParamTypeTime = 25		// 匹配时间格式
	ParamTypeMD5 = 26		// MD5类型字符串
	ParamTypeDiy = 100		// 只跑自定义的参数检查函数
)


// 验证参数是否合法
func checkParam(_ptype int, _val string) string {

	var pass = false

	switch _ptype {
	case ParamTypeAll:
		pass = true
	case ParamTypeSafe:
		pass = clCommon.CheckIsSafe(_val)
	case ParamTypeInt:
		pass = clCommon.CheckIsInt(_val)
	case ParamTypeFloat:
		pass = clCommon.CheckIsFloat(_val)
	case ParamTypeHtml:				// 通过htmlSpecialChars编码
		_val = clCommon.HtmlSpecialChars(_val)
		pass = true
	case ParamTypeUrl:				// 通过URL编码
		_val = url.QueryEscape(_val)
		pass = true
	case ParamTypePhone:
		pass = clCommon.CheckPhone(_val)
	case ParamTypeEmail:
		pass = clCommon.CheckEmail(_val)
	case ParamTypeDomain:
		pass = clCommon.CheckDomain(_val)
	case ParamTypeIP:
		pass = clCommon.CheckIP(_val)
	case ParamTypeBankCard:
		pass = clCommon.CheckBankCard(_val)
	case ParamTypeDate:
		pass = clCommon.CheckDate(_val)
	case ParamTypeTime:
		pass = clCommon.CheckTime(_val)
	case ParamTypeDateTime:
		pass = clCommon.CheckDateTime(_val)
	case ParamTypeUserPass:
		pass = clCommon.CheckUsername(_val)
	case ParamTypeQQ:
		pass = clCommon.CheckQQ(_val)
	case ParamTypeWechat:
		pass = clCommon.CheckWechat(_val)
	case ParamTypeImageName:
		pass = clCommon.CheckImageName(_val)
	case ParamTypeVcode:
		pass = clCommon.CheckIsVcode(_val)
	case ParamTypeClientType:
		pass = clCommon.CheckIsClient(_val)
	case ParamTypeUUID:
		pass = clCommon.CheckIsUUID(_val)
	case ParamTypeVersion:
		pass = clCommon.CheckIsVersion(_val)
	case ParamTypeImageBase64:
		pass = clCommon.CheckIsImageBase64(_val)
	case ParamTypeMD5:
		pass = clCommon.CheckIsMD5(_val)
	case ParamTypeDiy:
		pass = true
	}

	if pass {
		return _val
	} else {
		return "VALUE_IS_ERROR"
	}
}


