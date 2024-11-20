package BLC

import "fmt"

// 获取地址列表
func (cli *CLI) GetAccounts(nodeID string) {
	wallets := NewWallets(nodeID)
	fmt.Println("\t账号列表")
	for key, _ := range wallets.Wallets {
		fmt.Printf("\t\t[%s]\n", key)
	}
}
