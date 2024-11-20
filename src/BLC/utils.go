package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// 实现int64->[]byte
func Int64ToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int transact to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// 标准JSON格式转切片
// windows下需要添加引号
func JSON2Slice(jsonString string) []string {
	var strSlice []string
	if err := json.Unmarshal([]byte(jsonString), &strSlice); err != nil {
		log.Panicf("json to []string failed! %v\n", err)
	}
	return strSlice
}

// 参数数量检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		// 直接退出
		os.Exit(1)
	}
}

// string 2 hash160
func String2Hash160(address string) []byte {
	pubKeyHash := Base58Decode([]byte(address))

	return pubKeyHash[:len(pubKeyHash)-addressChecksumLen]
}

// 获取节点ID
func GetEnvNodeID() string {
	nodeId := os.Getenv("NODE_ID")
	if nodeId == "" {
		fmt.Println("NODE_ID env var not set...")
		os.Exit(1)
	}
	return nodeId
}

// gob编码
func gobEncode(data interface{}) []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// 命令转换为请求([]byte)
func commandToBytes(cmd string) []byte {
	var bytes [COMMAND_LENGTH]byte
	for i, c := range cmd {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

// 反解析，把请求中的命令解析出来
func bytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x00 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}
