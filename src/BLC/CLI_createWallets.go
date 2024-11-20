package BLC

import "fmt"

// 创建钱包集合
func (cli *CLI) CreateWallets(nodeID string) {
	wallets := NewWallets(nodeID)
	wallets.CreateWallet(nodeID)
	fmt.Printf("wallets : %v\n", wallets)
}
