package BLC

// 初始化区块链
func (cli *CLI) createBlockChain(address string) {
	CreateBlockChainWithGenesisBlock(address)
}
