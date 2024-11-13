package BLC

// UTXO结构管理
type UTXO struct {
	TxHash []byte    //UTXO对应的交易hash
	Index  int       //UTXO在其所属交易的输出列表中的索引
	Output *TxOutput //Output本身
}
