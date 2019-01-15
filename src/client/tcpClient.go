package client

import (
	"net"
	"fmt"
	"encoding/json"
	"time"
	"math/rand"
	"bufio"
	"base"
	"protocol/tcp"
)

type TcpFileSync struct {
	FileSyncer *base.FileSync
	IsSync     bool
	TcpClient
	tcp.PkgHandle
}

// tcp 客户端连接
type TcpClient struct {
	Connection		*net.TCPConn
	HawkServer		*net.TCPAddr
}

func run(fileSyncer *base.FileSync) (*TcpFileSync,error) {

	tcpAddr ,err := net.ResolveTCPAddr("tcp", fileSyncer.ServerAddr)

	if err != nil {
		fmt.Printf("hawk server[%s] resolve error [%s]\n", fileSyncer.ServerAddr,err.Error())
		return nil,err
	}

	conn, err := net.DialTCP("tcp",nil,tcpAddr)
	if err != nil {
		fmt.Printf("connect to hawk server [%s] error [%s]!\n", fileSyncer.ServerAddr,err.Error())
		return nil,err
	}

	client := TcpClient{
		Connection:conn,
		HawkServer:tcpAddr,
	}

	tcpFileSync := &TcpFileSync{
		TcpClient : client,
		FileSyncer: fileSyncer,
	}

	// 完成创建之后，启动单独的协程处理来自服务端的数据处理
	go tcpFileSync.dataReader()

	return tcpFileSync,nil
}

func (tcpFileSyncer *TcpFileSync) dataReader() {
	defer tcpFileSyncer.Connection.Close()
	bufferReader := bufio.NewReader(tcpFileSyncer.Connection)
	tcpFileSyncer.Unpack(bufferReader)
}


// 发送数据包底层原型
func (tcpFileSyncer *TcpFileSync) sender(pkgType byte,pkgContent []byte) error  {

	packet := &tcp.Packet{
		PacketType:pkgType,
		PacketContent:pkgContent,
	}

	pkgData,err := json.Marshal(packet)

	if err != nil {
		fmt.Printf("json marshal err [%s]\n",err.Error())
		return err
	}
	
	_,err = tcpFileSyncer.TcpClient.Connection.Write(tcpFileSyncer.Packaged(pkgData))

	return err
}

// 发送一个心跳包
func (tcpConsumer *TcpFileSync) SendHeart() ([]byte,error)  {

	heartPacket := &tcp.HeartPacket{
		Version: "1.0",
		Timestamp:time.Now().Unix(),
	}

	heartPacketJson,err := json.Marshal(heartPacket)

	if err != nil {
		fmt.Println("json marshal err")
		return nil,err
	}

	return nil,tcpConsumer.sender(base.HEART_BEAT_PACKET,heartPacketJson)
}

// 发送一个日志消息数据包
func (tcpConsumer *TcpFileSync) SendMessage(message string) error  {

	messagePacket := &tcp.MessagePacket{
		Content:message,
		Rand:rand.Int(),
		Timestamp:time.Now().Unix(),
	}

	messagePacketJson,err := json.Marshal(messagePacket)

	if err != nil {
		fmt.Println("json marshal err")
		return err
	}

	return tcpConsumer.sender(base.MESSAGE_PACKET,messagePacketJson)
}


/**
	发送心跳数据包
*/
func (tcpFileSyncer *TcpFileSync)Heart()  {
	for {
		tcpFileSyncer.SendHeart()
		time.Sleep(time.Second * 5)
	}
}


func (tcpFileSyncer *TcpFileSync)Dispatch(recvBuffer []byte) {
	var packet tcp.Packet
	if json.Unmarshal(recvBuffer,&packet) != nil {
		fmt.Println("unmarshal err")
	}
	switch packet.PacketType {
	case base.HEART_BEAT_PACKET:
		var heartPacket tcp.HeartPacket
		if json.Unmarshal(packet.PacketContent,&heartPacket) != nil {
			fmt.Println("client:json unmarshal error during HEART_BEAT_PACKET!")
			return
		}
		lastHeart := <-tcpFileSyncer.FileSyncer.HeartTime
		now := time.Now().Unix()
		if heartPacket.Timestamp > lastHeart && now - heartPacket.Timestamp > 10 {
			fmt.Printf("false: now [%d]  last [%d] \n",now,lastHeart)
			tcpFileSyncer.IsSync = false
		} else {
			fmt.Println("true")
			tcpFileSyncer.IsSync = true
		}
		tcpFileSyncer.FileSyncer.HeartTime <- now
		break
	case base.MESSAGE_PACKET:
		// 普通消息暂时不做处理
		break
	}

}
