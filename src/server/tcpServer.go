package server

import (
	"base"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"protocol/tcp"
	"time"
)

type TcpFileSync struct {
	FileSyncer *base.FileSync
	TcpServer
	tcp.PkgHandle
	conn net.Conn
}

// tcp 服务端连接
type TcpServer struct {
	listener   *net.TCPListener
	hawkServer *net.TCPAddr
}

func Run(fileSyncer *base.FileSync) {
	hawkServer, err := net.ResolveTCPAddr("tcp", fileSyncer.ServerAddr)
	if err != nil {
		fmt.Printf("hawk server[%s] resolve error [%s] \n", fileSyncer.ServerAddr, err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp", hawkServer)
	if err != nil {
		fmt.Printf("set server listen error，server [%s] err [%s] \n", fileSyncer.ServerAddr, err.Error())
		return
	}

	defer listener.Close()

	tcpServer := TcpServer{
		listener:   listener,
		hawkServer: hawkServer,
	}

	fmt.Printf("tcp server has started!serverAddr: %s ,waiting clients...\n", fileSyncer.ServerAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("a connection has been interrupted!")
			continue
		}

		tcpFileSync := &TcpFileSync{
			TcpServer:  tcpServer,
			FileSyncer: fileSyncer,
			conn:       conn,
		}

		tcpFileSync.PkgHandle.Dispatch = tcpFileSync.Dispatch

		// 服务端与客户端成功建立连接之后，启动新协程单独处理 Tcp 数据通信
		go tcpFileSync.dataReader()
	}
}

func (tcpFileSync *TcpFileSync) dataReader() {
	defer tcpFileSync.conn.Close()
	bufferReader := bufio.NewReader(tcpFileSync.conn)
	tcpFileSync.Unpack(bufferReader)
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
			fmt.Println("json unmarshal error during HEART_BEAT_PACKET!")
			return
		}
		lastHeart := <-tcpFileSync.FileSyncer.HeartTime
		now := time.Now().Unix()
		if heartPacket.Timestamp > lastHeart && now-heartPacket.Timestamp > 10 {
			fmt.Printf("false: now [%d]  last [%d] \n", now, lastHeart)
			tcpFileSync.FileSyncer.IsSync = false
		} else {
			fmt.Println("true")
			tcpFileSync.FileSyncer.IsSync = true
		}
		tcpFileSync.FileSyncer.HeartTime <- heartPacket.Timestamp

		// 返回一个全新的心跳包，让客户端处理即可
		heartPacket.Timestamp = time.Now().Unix()
		retBytes, err := json.Marshal(heartPacket)
		checkErr(err)

		packet.PacketContent = retBytes
		err = tcpFileSync.sender(packet)
		checkErr(err)
		break
	case base.MESSAGE_PACKET:
		var messagePacket tcp.MessagePacket
		if json.Unmarshal(packet.PacketContent, &messagePacket) != nil {
			fmt.Println("json unmarshal error during MESSAGE_PACKET!")
			return
		}
		content := messagePacket.Content
		fmt.Printf("server has received data [%s] \n", content)

		retData := &base.JsonResult{
			Status:  1,
			Message: "ok",
		}
		retBytes, err := json.Marshal(retData)
		checkErr(err)

		packet.PacketContent = retBytes
		err = tcpFileSync.sender(packet)
		checkErr(err)
		break
	}

}

// 发送数据包底层原型
func (tcpFileSync *TcpFileSync) sender(packet tcp.Packet) error {
	pkgData, err := json.Marshal(packet)

	if err != nil {
		fmt.Printf("json marshal err [%s]\n", err.Error())
		return err
	}

	_, err = tcpFileSync.conn.Write(tcpFileSync.Packaged(pkgData))

	return err
}

// check err
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
