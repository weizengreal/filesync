package client

import (
	"base"
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"protocol/tcp"
	"time"
)

type TcpFileSync struct {
	FileSyncer *base.FileSync
	TcpClient
	tcp.PkgHandle
}

// tcp 客户端连接
type TcpClient struct {
	Connection *net.TCPConn
	HawkServer *net.TCPAddr
}

func RunTcp(fileSyncer *base.FileSync) (*TcpFileSync, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", fileSyncer.ServerAddr)

	if err != nil {
		fmt.Printf("hawk server[%s] resolve error [%s]\n", fileSyncer.ServerAddr, err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("connect to hawk server [%s] error [%s]!\n", fileSyncer.ServerAddr, err.Error())
		return nil, err
	}

	fmt.Println("localAddr:", conn.LocalAddr())

	client := TcpClient{
		Connection: conn,
		HawkServer: tcpAddr,
	}

	tcpFileSync := &TcpFileSync{
		TcpClient:  client,
		FileSyncer: fileSyncer,
	}

	tcpFileSync.PkgHandle.Dispatch = tcpFileSync.Dispatch

	// 完成创建之后，启动单独的协程处理来自服务端的数据处理
	go tcpFileSync.dataReader()

	// 心跳
	go tcpFileSync.Heart()

	return tcpFileSync, nil
}

func (tcpFileSync *TcpFileSync) HangTcpClient() {
	// loop forever
	for {
		select {
		case <-tcpFileSync.FileSyncer.StopChan:
			fmt.Println("receive stop signal,return main func!")
			return
		}

		// 先放这儿
		time.Sleep(1 * time.Second)
	}
}

func (tcpFileSync *TcpFileSync) dataReader() {
	defer tcpFileSync.Connection.Close()
	bufferReader := bufio.NewReader(tcpFileSync.Connection)
	tcpFileSync.Unpack(bufferReader)
	tcpFileSync.FileSyncer.StopChan <- struct{}{}
	fmt.Println("client finished!push to stop chan!")
}

// 发送数据包底层原型
func (tcpFileSync *TcpFileSync) sender(pkgType byte, pkgContent []byte) error {

	packet := &tcp.Packet{
		PacketType:    pkgType,
		PacketContent: pkgContent,
	}

	pkgData, err := json.Marshal(packet)

	if err != nil {
		fmt.Printf("json marshal err [%s]\n", err.Error())
		return err
	}

	_, err = tcpFileSync.TcpClient.Connection.Write(tcpFileSync.Packaged(pkgData))

	return err
}

// 发送一个心跳包
func (tcpConsumer *TcpFileSync) SendHeart() error {

	heartPacket := &tcp.HeartPacket{
		Version:   "1.0",
		Timestamp: time.Now().Unix(),
	}

	heartPacketJson, err := json.Marshal(heartPacket)

	if err != nil {
		fmt.Println("json marshal err")
		return err
	}

	return tcpConsumer.sender(base.HEART_BEAT_PACKET, heartPacketJson)
}

// 发送一个日志消息数据包
func (tcpConsumer *TcpFileSync) SendMessage(message string) error {

	messagePacket := &tcp.MessagePacket{
		Content:   message,
		Rand:      rand.Int(),
		Timestamp: time.Now().Unix(),
	}

	messagePacketJson, err := json.Marshal(messagePacket)

	if err != nil {
		fmt.Println("json marshal err")
		return err
	}

	return tcpConsumer.sender(base.MESSAGE_PACKET, messagePacketJson)
}

/**
TODO:: tcp 改造为队列读取，可以在繁忙时省下一个心跳包的消耗
发送心跳数据包
*/
func (tcpFileSync *TcpFileSync) Heart() {
	for {
		tcpFileSync.SendHeart()
		time.Sleep(time.Second * 5)
	}
}

func (tcpFileSync *TcpFileSync) Dispatch(recvBuffer []byte) {
	var packet tcp.Packet
	if json.Unmarshal(recvBuffer, &packet) != nil {
		fmt.Println("unmarshal err")
	}
	switch packet.PacketType {
	case base.HEART_BEAT_PACKET:
		var heartPacket tcp.HeartPacket
		if json.Unmarshal(packet.PacketContent, &heartPacket) != nil {
			fmt.Println("client:json unmarshal error during HEART_BEAT_PACKET!")
			return
		}
		lastHeart := <-tcpFileSync.FileSyncer.HeartTime
		now := time.Now().Unix()
		if heartPacket.Timestamp > lastHeart && now-heartPacket.Timestamp > 10 {
			fmt.Printf("false: now [%d]  last [%d] \n", now, lastHeart)
			tcpFileSync.FileSyncer.IsSync = false
		} else {
			tcpFileSync.FileSyncer.IsSync = true
		}
		tcpFileSync.FileSyncer.HeartTime <- heartPacket.Timestamp
		break
	case base.MESSAGE_PACKET:
		// 普通消息暂时不做处理
		break
	}

}
