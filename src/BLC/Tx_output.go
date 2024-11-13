package BLC

// 交易的输出管理
type TxOutput struct {
	Value        int    //金额
	ScriptPubkey string //用户名（UTXO所有者）
}

// 验证当前UTXO是否数以指定的地址
func (txOutput *TxOutput) CheckPubkeyWithAddress(address string) bool {
	return address == txOutput.ScriptPubkey
}
