package BLC

import (
	"fmt"
	"os"
)

// 发起交易函数
func (cli *CLI) send(from []string, to []string, amount []string) {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	//获取区块链对象
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	if len(from) != len(to) || len(from) != len(amount) {
		fmt.Println("交易参数输入有误，请检查一致性...")
		os.Exit(1)
	}
	blockchain.MineNewBlock(from, to, amount)
}
