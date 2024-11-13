package BLC

import (
	"fmt"
	"os"
)

// 打印完整区块链信息
func (cli *CLI) printchain() {
	if !dbExist() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	blockchain.PrintChain()
}
