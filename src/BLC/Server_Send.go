package BLC

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

//请求发送文件

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
	_, err = io.Copy(conn, bytes.NewReader(msg))
	if err != nil {
		log.Panicf("add data to conn  failed: %v\n", err)
	}
}

// 区块链版本验证
func sendVersion(toAddress string, bc *BlockChain) {
	// 1. 获取当前节点区块高度
	height := bc.GetHeight()
	// 2.组装生成version
	versionData := Version{int(height), nodeAddress}

	// 3. 组装要发送的请求
	data := gobEncode(versionData)
	// 4. 将命令与版本组装成完整的请求
	request := append(commandToBytes(CMD_VERSION), data...)
	// 4. 发送请求
	SendMessage(toAddress, nodeAddress, request)
}

// 从指定节点同步数据
func sendGetBlocks(toAddress string) {
	// 1. 生成数据
	data := gobEncode(GetBlocks{AddressFrom: nodeAddress})
	// 2. 组装请求
	request := append(commandToBytes(CMD_GETBLOCKS), data...)
	// 3. 发送请求
	SendMessage(toAddress, nodeAddress, request)
}

// 发送指定区块请求
func sendGetData(toAddress string, hash []byte) {
	// 1. 生成数据
	data := gobEncode(GetData{AddressFrom: nodeAddress, ID: hash})
	// 2. 组装请求
	request := append(commandToBytes(CMD_GETDATA), data...)
	// 3. 发送请求
	SendMessage(toAddress, nodeAddress, request)

}

// 发送区块展示，向其他节点展示
func sendInv(toAddress string, hashes [][]byte) {
	// 1. 生成数据
	data := gobEncode(Inv{AddressFrom: nodeAddress, Hashes: hashes})
	// 2. 组装请求
	request := append(commandToBytes(CMD_INV), data...)
	// 3. 发送请求
	SendMessage(toAddress, nodeAddress, request)

}

// 发送区块信息
func sendBlock(toAddress string, block []byte) {
	// 1. 生成数据
	data := gobEncode(BlockData{AddressFrom: nodeAddress, Block: block})
	// 2. 组装请求
	request := append(commandToBytes(CMD_BLOCK), data...)
	// 3. 发送请求
	SendMessage(toAddress, nodeAddress, request)

}
