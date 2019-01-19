package client

import (
	"base"
	"fmt"
	"net"
	"os"
	"protocol/tcp"
)

type UdpFileSync struct {
	FileSyncer *base.FileSync
	conn       net.Conn
	tcp.PkgHandle
}

func RunUdp(fileSyncer *base.FileSync) {
	conn, err := net.Dial("udp", fileSyncer.ServerAddr)
	defer conn.Close()
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("localAddr:", conn.LocalAddr())

	conn.Write([]byte("Hello world!"))

	fmt.Println("send msg")

	var msg [20]byte
	conn.Read(msg[0:])

	fmt.Println("msg is:", string(msg[0:10]))

	//udpFileSync := &UdpFileSync{
	//	FileSyncer: fileSyncer,
	//	conn:       conn,
	//}
	//udpFileSync.PkgHandle.Dispatch = udpFileSync.Dispatch
	//
	//return udpFileSync, nil
}

//
//
//func (udpFileSync *UdpFileSync) dataReader() {
//	defer udpFileSync.conn.Close()
//	bufferReader := bufio.NewReader(udpFileSync.conn)
//	udpFileSync.Unpack(bufferReader)
//	udpFileSync.FileSyncer.StopChan <- struct{}{}
//	fmt.Println("client finished!push to stop chan!")
//}
//
//
//func (udpFileSync *UdpFileSync) HangUdpClient() {
//	// loop forever
//	for {
//		select {
//		case <-udpFileSync.FileSyncer.StopChan:
//			fmt.Println("receive stop signal,return main func!")
//			return
//		}
//	}
//}
//
//
//func (udpFileSync *UdpFileSync) Dispatch(recvBuffer []byte) {
//	fmt.Println("receive msg:",string(recvBuffer))
//}
