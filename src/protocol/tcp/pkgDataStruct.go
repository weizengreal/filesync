package tcp


// 数据包格式
type Packet struct {
	PacketType		byte `json:"packetType"`
	PacketContent	[]byte `json:"packetContent"`
}

// 心跳包
type HeartPacket struct {
	Version			string `json:"version"`
	Timestamp		int64 `json:"timestamp"`
}

// 消息体包
type MessagePacket struct {
	Content			string `json:"content"`
	Rand 			int `json:"rand"`
	Timestamp		int64 `json:"timestamp"`
}


