package main

import (
	"clws-framework/core/clGlobal"
	"clws-framework/core/clWebsocket"
	"clws-framework/rule"
	"fmt"
)

func main() {

	rule.InitRule()

	clGlobal.Init("sky.conf")
	err := clWebsocket.Serve("websocket", 16666)

	fmt.Printf("服务器退出: %v\n", err)
}
