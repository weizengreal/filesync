package server

import (
	"base"
	"fmt"
	"net"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func recvUDPMsg(conn *net.UDPConn) {
	var buf [20]byte

	n, raddr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	fmt.Println("msg is ", string(buf[0:n]))

	_, err = conn.WriteToUDP([]byte("nice to see u"), raddr)
	checkError(err)
}

func RunUdp(fileSyncer *base.FileSync) {

	udp_addr, err := net.ResolveUDPAddr("udp", fileSyncer.ServerAddr)
	checkError(err)

	conn, err := net.ListenUDP("udp", udp_addr)
	defer conn.Close()
	checkError(err)

	recvUDPMsg(conn)
}
