package BLC

// 初始化区块链
func (cli *CLI) createBlockChain(address string) {
	bc := CreateBlockChainWithGenesisBlock(address)
	defer bc.DB.Close()
	//设置utxo重置操作
	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}
