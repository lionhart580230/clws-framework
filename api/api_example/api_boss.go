package api_example

import (
	"clws-framework/core/clDebug"
	"clws-framework/core/clPacket"
	"clws-framework/core/clRouter"
	"clws-framework/core/clUserPool"
)

// Api例子
func ApiExample(_uInfo *clUserPool.ClNetUserInfo, _params map[string]string) *clPacket.RuleCBResp {
	clDebug.Debug("参数列表: %+v", _params)
	return clRouter.JCode("apiExampleResp", "", nil)
}