package BLC

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//网络服务文件管理

//3000作为引导节点（主节点地址）
var knowNodes = []string{"localhost:3000"}

//节点地址
var nodeAddress string

//启动服务
func stratServer(nodeID string) {
	fmt.Printf("----------Start Node[%v]----------\n", nodeID)
	//节点地址赋值
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	//1.监听节点
	listen, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil {
		log.Panicf("listen address of %s failed: %v\n!", nodeAddress, err)
	}
	defer listen.Close()
	// 两个节点，主节点负责保存数据，钱包节点负责发送请求，同步数据
	if nodeAddress != knowNodes[0] {
		// 1. 不是主节点，发送请求，同步数据
		// ...
		//SentMessage(knowNodes[0], nodeAddress)
		sendVersion(knowNodes[0])
	}

	for {
		// 2. 生成链接，接受请求
		conn, err := listen.Accept()
		if err != nil {
			log.Panicf("accept connect failed: %v", err)
		}
		//处理请求
		// 单独启动一个goroutine 进行请求处理
		go handleConnection(conn)
	}
}

// worker
// 请求处理函数
func handleConnection(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panicf("Receive a request failed: %v\n", err)
	}
	cmd := bytesToCommand(request[:COMMAND_LENGTH])
	fmt.Printf("Receive a command %s\n!", cmd)
	switch cmd {
	case CMD_VERSION:
		handleVersion()
	case CMD_GETDATA:
		handleGetData()
	case CMD_GETBLOCKS:
		handleGetBlocks()
	case CMD_INV:
		handleInv()
	case CMD_BLOCK:
		handleBlock()
	default:
		fmt.Printf("Command not recognized\n")

	}
}
