package BLC

import (
	"bytes"
	"fmt"
	"io"
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
		request, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panicf("read connect request failed: %v", err)
		}
		// 3. 处理请求
		fmt.Printf("Receive a message %v\n", request)
		handleConnection()
	}
}

// 请求处理函数
func handleConnection() {

}

//发送消息
func SendMessage(to, from string, msg []byte) {
	fmt.Printf("Sending request %v to %v\n", from, to)
	// 1.连接上服务器
	conn, err := net.Dial(PROTOCOL, to)
	if err != nil {
		log.Panicf("connect to %v failed: %v\n", to, err)
	}
	defer conn.Close()
	// 要发送的数据
	_, err = io.Copy(conn, bytes.NewReader([]byte(msg)))
	if err != nil {
		log.Panicf("add data to conn  failed: %v\n", err)
	}

}

// 区块链版本验证
func sendVersion(toAddress string) {
	// 1. 获取当前节点区块高度
	height := 1
	// 2.组装生成version
	versionData := Version{height, nodeAddress}

	// 3. 组装要发送的请求
	data := gobEncode(versionData)
	// 4. 将命令与版本组装成完整的请求
	request := append(commandToBytes(CMD_VERSION), data...)
	// 4. 发送请求
	SendMessage(toAddress, nodeAddress, request)
}
