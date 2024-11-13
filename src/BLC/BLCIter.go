package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// 区块链迭代管理文件

// BlockChainIter 迭代器基本结构
type BlockChainIter struct {
	DB          *bolt.DB //迭代目标
	CurrentHash []byte   // 当前迭代目标Hash

}

// Iterator 创建迭代器对象
func (blc *BlockChain) Iterator() *BlockChainIter {
	return &BlockChainIter{blc.DB, blc.Tip}
}

// Next 实现迭代函数next，获取到每一个区块
func (bcit *BlockChainIter) Next() *Block {
	var block *Block

	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBlockBytes := b.Get(bcit.CurrentHash)
			block = Deserialize(currentBlockBytes)
			// 更新迭代器中区块哈希
			bcit.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panicf("iterato the db failed %v\n", err)
	}
	return block
}
