package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//交易管理文件

// 定义一个基本结构
type Transaction struct {
	TxHash []byte      //交易哈希（标识）
	Vins   []*TxINPUT  //输入的列表
	Vouts  []*TxOutput //输出的列表

}

// 实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	//输入
	//coinbase特点
	//txHash：nil
	//vout：
	//ScriptSig
	txInput := &TxINPUT{[]byte{}, -1, "system reward"}
	//输出：
	//value：
	//address:
	txOutput := &TxOutput{10, address}
	txCoinbase := &Transaction{nil, []*TxINPUT{txInput}, []*TxOutput{txOutput}}
	txCoinbase.HashTransaction()
	return txCoinbase
}

// 生成交易哈希（交易序列化）
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panic("tx Hash failed %V\n", err)
	}
	//生成哈希值
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// 生成一个普通转账交易
func NewSimpleTransaction(from string, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxINPUT
	var txOutputs []*TxOutput

	// 调用可花费UTXO函数
	money, spendableUTXODic := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("money: %v\n", money)
	//输入
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panicf("decode string to []byte failed %V\n", err)
		}
		// 遍历索引列表
		for _, index := range indexArray {
			txInput := &TxINPUT{txHashBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}
	//输出(转账源)
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)
	//输出（找零）
	if money > amount {
		txOutput = &TxOutput{money - amount, from}
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足...\n")
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction()
	return &tx
}

// 判断一个指定的交易是否时一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}
