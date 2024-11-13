package BLC

import "fmt"

// 创建钱包集合
func (cli *CLI) CreateWallets() {
	wallets := NewWallets()
	wallets.CreateWallet()
	fmt.Printf("wallets : %v!\n", wallets)
}
