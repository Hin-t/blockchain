package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//UTXO持久化管理

//用于存放UTXO的bucket
const utxoTableName = "utxoTable"

//utxo_set结构（保存指定区块链中所有的UTXO）

type UTXOSet struct {
	BlockChain *BlockChain
}

//更新

//查找
func ()
//重置
func (utxoSet *UTXOSet) ResetUTXOSet() {
	//在第一次创建的时候就更新utxo table
	utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		// 查找utxo table
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := b.Delete([]byte(utxoTableName))
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
