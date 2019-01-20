package test

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"testing"
)

func TestReusePort(t *testing.T) {
	localAddr := "127.0.0.1:11198"

	conn1, err := reuseport.Dial("udp", localAddr, "115.159.225.216:11197")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer conn1.Close()

	conn2, err := reuseport.Dial("udp", localAddr, "134.175.14.99:11197")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer conn2.Close()

	conn1.Write([]byte("Hello world!"))

	conn2.Write([]byte("Hello world!"))

}
