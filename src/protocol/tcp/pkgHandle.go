package tcp

import (
	"bufio"
	"io"
	"fmt"
	"hash/crc32"
)

/**
	定义数据包的各种操作
 */
type pkgData interface {

	Packaged(sendBytes []byte) []byte

	Unpack(bufferReader *bufio.Reader)

	Dispatch(recvBuffer []byte)
}

/**
	定义数据包的操作结构体，供外部实例调用
 */
type PkgHandle struct {
	pkgData
}

/**
	根据传递进来的 reader 读取字节流
 */
func (pkgHandle *PkgHandle)Unpack(bufferReader *bufio.Reader) {
	// 状态机状态
	state := 0x00
	// 数据包长度
	length := uint16(0)
	// crc校验和
	crc16 := uint16(0)
	var recvBuffer []byte
	// 游标
	cursor := uint16(0)
	//状态机处理数据
	for {
		recvByte,err := bufferReader.ReadByte()
		if err != nil {
			//这里因为做了心跳，所以就没有加deadline时间，如果远程端断开连接
			//这里ReadByte方法返回一个io.EOF的错误，具体可考虑文档
			if err == io.EOF {
				fmt.Println("get EOF string!")
			}
			//在这里直接退出goroutine，关闭连接由defer操作完成
			return
		}
		//进入状态机，根据不同的状态来处理
		switch state {
		case 0x00:
			if recvByte == 0xFF {
				state = 0x01
				//初始化状态机
				recvBuffer = nil
				length = 0
				crc16 = 0
			}else{
				state = 0x00
			}
			break
		case 0x01:
			if recvByte == 0xFF {
				state = 0x02
			}else{
				state = 0x00
			}
			break
		case 0x02:
			length += uint16(recvByte) * 256
			state = 0x03
			break
		case 0x03:
			length += uint16(recvByte)
			// 一次申请缓存，初始化游标，准备读数据
			recvBuffer = make([]byte,length)
			cursor = 0
			state = 0x04
			break
		case 0x04:
			//不断地在这个状态下读数据，直到满足长度为止
			recvBuffer[cursor] = recvByte
			cursor++
			if(cursor == length){
				state = 0x05
			}
			break
		case 0x05:
			crc16 += uint16(recvByte) * 256
			state = 0x06
			break
		case 0x06:
			crc16 += uint16(recvByte)
			state = 0x07
			break
		case 0x07:
			if recvByte == 0xFF {
				state = 0x08
			}else{
				state = 0x00
			}
		case 0x08:
			if recvByte == 0xFE {
				//执行数据包校验
				if (crc32.ChecksumIEEE(recvBuffer) >> 16) & 0xFFFF == uint32(crc16) {
					pkgHandle.Dispatch(recvBuffer)
				}else{
					fmt.Println("drop this data!receive:",string(recvBuffer))
				}
			}
			//状态机归位,接收下一个包
			state = 0x00
		}
	}
}


/**
	数据包封装函数
 */
func (pkgHandle *PkgHandle)Packaged(sendBytes []byte) []byte {
	packetLength := len(sendBytes) + 8
	result := make([]byte,packetLength)
	result[0] = 0xFF
	result[1] = 0xFF
	result[2] = byte(uint16(len(sendBytes)) >> 8)
	result[3] = byte(uint16(len(sendBytes)) & 0xFF)
	copy(result[4:],sendBytes)
	sendCrc := crc32.ChecksumIEEE(sendBytes)
	result[packetLength-4] = byte(sendCrc >> 24)
	result[packetLength-3] = byte(sendCrc >> 16 & 0xFF)
	result[packetLength-2] = 0xFF
	result[packetLength-1] = 0xFE
	//fmt.Println(result)
	return result
}
