package main

import (
	"fmt"
	"github.com/xiaolan580230/clws-framework/core/clGlobal"
	"github.com/xiaolan580230/clws-framework/core/clWebsocket"
	"github.com/xiaolan580230/clws-framework/rule"
)

func main() {

	rule.InitRule()

	clGlobal.Init("sky.conf")
	err := clWebsocket.Serve("websocket", 16666)

	fmt.Printf("服务器退出: %v\n", err)
}
