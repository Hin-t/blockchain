package BLC

import "bytes"

// 交易的输出管理
type TxOutput struct {
	Value int //金额
	//ScriptPubkey string //用户名（UTXO所有者）
	//用户名（UTXO所有者）
	Ripemd160Hash []byte
}

// 验证当前UTXO是否数以指定的地址
//func (txOutput *TxOutput) CheckPubkeyWithAddress(address string) bool {
//	return address == txOutput.ScriptPubkey
//}

// output身份验证
func (txOutput *TxOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	//转换
	hash160 := String2Hash160(address)
	return bytes.Compare(hash160, txOutput.Ripemd160Hash) == 0
}

// 新建output对象
func NewTxOutput(value int, address string) *TxOutput {
	txOutput := &TxOutput{value, String2Hash160(address)}
	return txOutput
}
