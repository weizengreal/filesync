package base

/*
	数据发送者接口

	这里定义了数据发送者必须具有的基本方法，发送普通数据和心跳数据
*/
type Sender interface {
	// 发送简单的数据包
	SendMessage(logs string) error

	// 发送心跳数据包
	SendHeart() ([]byte,error)
}

/*
	文件同步基本信息

	这里定义了数据的消费者必须拥有的基本字段
*/
type FileSync struct {
	ServerAddr 	string
	HeartTime	chan int64
	StopChan	chan struct{}

	Sender
}
