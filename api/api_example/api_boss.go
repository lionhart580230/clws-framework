package api_example

import (
	"github.com/xiaolan580230/clws-framework/core/clDebug"
	"github.com/xiaolan580230/clws-framework/core/clPacket"
	"github.com/xiaolan580230/clws-framework/core/clRouter"
	"github.com/xiaolan580230/clws-framework/core/clUserPool"
)

// Api例子
func ApiExample(_uInfo *clUserPool.ClNetUserInfo, _params map[string]string) *clPacket.RuleCBResp {
	clDebug.Debug("参数列表: %+v", _params)
	return clRouter.JCode("apiExampleResp", "", nil)
}