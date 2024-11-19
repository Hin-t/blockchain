package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//UTXO持久化管理

//用于存放UTXO的bucket
const utxoTableName = "utxoTable.db"

//utxo_set结构（保存指定区块链中所有的UTXO）

type UTXOSet struct {
	BlockChain *BlockChain
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

//输出集合反序列化
func DeserializeTXOutputs(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	if err := decoder.Decode(&txOutputs); err != nil {
		log.Panicf("deserialize the struct utxo failed! %v\n", err)
	}
	return &txOutputs
}

//查询余额
func (utxoSet *UTXOSet) GetBalance(address string) int {
	UTXOS := utxoSet.FindUTXOWithAddress(address)
	var amount int
	for _, utxo := range UTXOS {
		fmt.Printf("utxo-txhash:%x\n", utxo.TxHash)
		fmt.Printf("utxo-index:%x\n", utxo.Index)
		fmt.Printf("utxo-Ripemd160Hash:%x\n", utxo.Output.Ripemd160Hash)
		fmt.Printf("utxo-value:%x\n", utxo.Output.Value)
		amount += utxo.Output.Value
	}
	return amount
}

//查找
func (utxoSet *UTXOSet) FindUTXOWithAddress(address string) []*UTXO {
	var utxos []*UTXO
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		//1. 获取utxo table
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			//cursor -- 游标
			c := b.Cursor()
			// 通过游标遍历bolt数据库中的数据
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutputs(v)
				for _, utxo := range txOutputs.TXOutputs {
					if utxo.UnLockScriptPubKeyWithAddress(address) {
						utxp_single := UTXO{Output: utxo}
						utxos = append(utxos, &utxp_single)
					}
				}
			}

		}
		return nil
	})
	if err != nil {
		return nil
	}
	if nil != err {
		log.Panicf("find the utxo of [%s] failed! %v\n", address, err)
	}
	return utxos
}

//重置
func (utxoSet *UTXOSet) ResetUTXOSet() {
	//在第一次创建的时候就更新utxo table
	utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		// 查找utxo table
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panicf("delete utxo table failed! %v\n", err)
			}
		}

		// 创建
		bucket, err := tx.CreateBucket([]byte(utxoTableName))
		if err != nil {
			log.Panicf("create bucket failed! %v\n", err)
		}
		if bucket != nil {
			// 查找当前所有utxo
			txOutputMap := utxoSet.BlockChain.FindUTXOMap()
			for keyHash, outputs := range txOutputMap {
				//将所有UTXO存入
				txHash, _ := hex.DecodeString(keyHash)
				fmt.Printf("KeyHash: %x\n", txHash)
				//存入utxo table
				err = bucket.Put(txHash, outputs.Serialize())
				if err != nil {
					log.Panicf("put utxo to bucket failed! %v\n", err)
				}
			}
		}
		return nil
	})
}

//更新
func (utxoSet *UTXOSet) Update() {
	//1获取最新区块
	latest_block := utxoSet.BlockChain.Iterator().Next()
	utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			// 只需查找最新一个区块的交易列表，因为每上链一个区块
			// utxo table都更新一次，所以只需查找最近一个区块中的交易
			for _, tx := range latest_block.Txs {
				if !tx.IsCoinbaseTransaction() {
					// 2.将已经被当前这笔交易输入引用的utxo删除
					for _, vin := range tx.Vins {
						// 需要更新的输出
						updateOutputs := TXOutputs{}
						// 获取指定输入所引用的交易哈希的输出
						outputBytes := b.Get(vin.TxHash)
						outs := DeserializeTXOutputs(outputBytes)
						for outIdx, out := range outs.TXOutputs {
							if vin.Vout != outIdx {
								updateOutputs.TXOutputs = append(updateOutputs.TXOutputs, out)
							}
						}
						// 如果交易中没有UTXO，删除该交易
						if len(updateOutputs.TXOutputs) == 0 {
							b.Delete(vin.TxHash)
						} else {
							// 将更新之后的存入数据库
							b.Put(vin.TxHash, updateOutputs.Serialize())
						}
					}

				}
				// 获取当前区块中新生成交易输出
				newOutputs := TXOutputs{}
				newOutputs.TXOutputs = append(newOutputs.TXOutputs, tx.Vouts...)
				b.Put(tx.TxHash, newOutputs.Serialize())
			}
			// 1.将最新区块中的UTXO插入

		}
		return nil
	})
}
