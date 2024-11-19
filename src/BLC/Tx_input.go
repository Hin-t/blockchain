package BLC

import "bytes"

// 交易的输入管理
// 输入的结构
type TxINPUT struct {
	TxHash []byte //交易的哈希
	Vout   int    //引用上一笔交易的输出索引号
	// 公钥
	PublicKey []byte
	// 数字签名
	Signature []byte
}

// 验证引用的地址是否匹配
//func (txInput *TxINPUT) CheckPubkeyWithAddress(address string) bool {
//	return address == txInput.ScriptSig
//}

// 传递哈希160进行判断
func (txInput *TxINPUT) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	// 获取Input的hash160
	inputRipemd160Hash := Ripemd160Hash(txInput.PublicKey)
	return bytes.Compare(inputRipemd160Hash, ripemd160Hash) == 0
}
