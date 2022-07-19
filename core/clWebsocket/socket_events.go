package clWebsocket

import (
	"github.com/xiaolan580230/clws-framework/core/clUserPool"
	"net/http"
)

// 连线进入
var eventOnAccept func(info *clUserPool.ClNetUserInfo, _r *http.Request) bool = nil

// 连线断开
var eventOnClose func(info *clUserPool.ClNetUserInfo) = nil

// 执行规则之前
var eventBeforeRule func(info *clUserPool.ClNetUserInfo, _param map[string] string) = nil

// 执行规则之后
var eventAfterRule func(info *clUserPool.ClNetUserInfo, _resp string) = nil


// 设置连线进入事件回调
func OnAcceptCallback(_func func(info *clUserPool.ClNetUserInfo, _r *http.Request) bool) {
	eventOnAccept = _func
}


// 设置连线断开事件回调
func OnCloseCallback(_func func(info *clUserPool.ClNetUserInfo)) {
	eventOnClose = _func
}


// 设置执行规则之前回调
func OnBeforeRuleCallback(_func func(info *clUserPool.ClNetUserInfo, _param map[string] string)) {
	eventBeforeRule = _func
}


// 设置执行规则之后回调
func OnAfterRuleCallback(_func func(info *clUserPool.ClNetUserInfo, _resp string)) {
	eventAfterRule = _func
}