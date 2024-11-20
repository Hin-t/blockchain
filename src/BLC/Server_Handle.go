package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// 请求处理文件管理

// Version
func handleVersion(requests []byte, bc *BlockChain) {
	fmt.Println("the request of version handle...")
	var buff bytes.Buffer
	var data Version
	// 1. 解析请求
	dataBytes := requests[COMMAND_LENGTH:]
	// 2. 生成version结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); err != nil {
		log.Panicf("decode the Version struct failed!%v\n", err)
	}
	// 3.获取请求方的区块高度
	versionHeight := data.Height
	// 4.获取自身节点的区块高度
	height := bc.GetHeight()
	// 如果当前节点的区块高度大于versionHeight
	// 将当前节点版本信息发送给请求节点
	if height > int64(versionHeight) {
		sendVersion(data.AddressFrom, bc)
	} else if height < int64(versionHeight) {
		// 如果当前节点区块高度小于versionHeight
		// 向发送发发起同步数据的请求
		sendGetBlocks(data.AddressFrom)
	}
}

// Getblocks
// 数据同步请求处理
func handleGetBlocks(requests []byte, bc *BlockChain) {
	fmt.Println("the request of get blocks handle...")
	var buff bytes.Buffer
	var data GetBlocks
	// 1. 解析请求
	dataBytes := requests[COMMAND_LENGTH:]
	// 2. 生成GetBlocks结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); err != nil {
		log.Panicf("decode the GetBlocks struct failed!%v\n", err)
	}
	// 3.获取区块链的所有的区块哈希
	hashes := bc.GetBlockHashes()
	sendInv(data.AddressFrom, hashes)
}

// GetData
// 处理获取指定区块的请求
func handleGetData(requests []byte, bc *BlockChain) {
	fmt.Println("the request of get data handle...")
	var buff bytes.Buffer
	var data GetData
	// 1. 解析请求
	dataBytes := requests[COMMAND_LENGTH:]
	// 2. 生成GetData结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); err != nil {
		log.Panicf("decode the GetData struct failed!%v\n", err)
	}
	// 3. 通过传来的区块哈希，获取本地节点的区块
	blockByte := bc.GetBlock(data.ID)
	sendBlock(data.AddressFrom, blockByte)
}

// Inv
func handleInv(requests []byte, bc *BlockChain) {
	fmt.Println("the request of Inv handle...")
	var buff bytes.Buffer
	var data Inv
	// 1. 解析请求
	dataBytes := requests[COMMAND_LENGTH:]
	// 2. 生成Inv结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); err != nil {
		log.Panicf("decode the Inv struct failed!%v\n", err)
	}
	for _, hash := range data.Hashes {
		sendGetData(data.AddressFrom, hash)
	}
}

// Block
// 接收到新区块进行处理
func handleBlock(requests []byte, bc *BlockChain) {
	fmt.Println("the request of BlockData handle...")
	var buff bytes.Buffer
	var data BlockData
	// 1. 解析请求
	dataBytes := requests[COMMAND_LENGTH:]
	// 2. 生成BlockData结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); err != nil {
		log.Panicf("decode the BlockData struct failed!%v\n", err)
	}
	// 3. 将接收到的区块添加到区块链中
	blockByte := data.Block
	block := Deserialize(blockByte)
	bc.AddBlock(block)
	// 4. 更新UTXO
	utxoSet := UTXOSet{bc}
	utxoSet.Update()
}
