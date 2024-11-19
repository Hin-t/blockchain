package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
