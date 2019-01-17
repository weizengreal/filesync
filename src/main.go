package main

import (
	"base"
	"client"
	"flag"
	"server"
	"time"
)

/**
入口函数，根据传入的参数决定当前是服务端还是客户端
*/
func main() {

	serverAddr := flag.String("serverAddr", "0.0.0.0:11197", "the path of server address!")
	appType := flag.Int("appType", 2, "application type,1 means server,2 means client!")
	flag.Parse()

	fileSyncer := &base.FileSync{
		ServerAddr: *serverAddr,
		HeartTime:  make(chan int64, 1),
		StopChan:   make(chan struct{}),
	}
	fileSyncer.HeartTime <- time.Now().Unix()

	if *appType == 1 {
		// 启动服务端
		fileSyncer.ServerAddr = "0.0.0.0:11197"
		server.Run(fileSyncer)
	} else if *appType == 2 {
		// 启动客户端
		tcpFileSyncer, err := client.Run(fileSyncer)
		if err != nil {
			panic(err.Error())
		}
		tcpFileSyncer.HangTcpClient()
	}

}
