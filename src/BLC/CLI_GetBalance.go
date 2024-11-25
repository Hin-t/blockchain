package BLC

import "fmt"

// 查询余额
func (cli *CLI) getBalance(from string, nodeID string) {
	//查找该地址UTXO
	//获取区块链对象
	blockchain := BlockchainObject(nodeID)
	//关闭实例对象
	defer blockchain.DB.Close()
	//amount := blockchain.getBalance(from)
	//fmt.Printf("\t地址[%s]的余额：[%d]\n", from, amount)

	//查找该地址UTXO
	utxoSet := UTXOSet{blockchain}
	amount := utxoSet.GetBalance(from)

	fmt.Printf("\t地址[%s]的余额：[%d]\n", from, amount)
}
