package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

//存入所有输出的集合
type TXOutputs struct {
	TXOutputs []*TxOutput
}

//输出集合序列化
func (txOutputs *TXOutputs) Serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	if err := encoder.Encode(txOutputs); err != nil {
		log.Panicf("serialize the utxo failed! %v\n", err)
	}
	return buff.Bytes()
}
