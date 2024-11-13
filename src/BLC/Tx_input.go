package BLC

// 交易的输入管理
// 输入的结构
type TxINPUT struct {
	TxHash    []byte //交易的哈希
	Vout      int    //引用上一笔交易的输出索引号
	ScriptSig string //用户名
}

// 验证引用的地址是否匹配
func (txInput *TxINPUT) CheckPubkeyWithAddress(address string) bool {
	return address == txInput.ScriptSig
}
