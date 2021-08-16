package clWebsocket


// 写出缓冲区对象
type WriteObj struct {
	data string			// 需要发送的数据
	connId uint64		// 需要发送的连接Id, 为0则不发送用户
}