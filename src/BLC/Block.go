package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// 区块基本结构与功能管理文件
//
// 实现一个最基本的区块结构

type Block struct {
	TimeStamp     int64  // 区块时间戳，代表区块时间
	Hash          []byte // 当前区块hash
	PrevBlockHash []byte // 前一个区块hash
	Height        int64  // 区块高度
	// Data          []byte         // 交易数据
	MerkleRoot []byte         // Merkle Root
	Txs        []*Transaction // 交易数据（交易列表）
	Nonce      int64          // 运行pow生成hash的变化值，也代表pow运行时动态修改的数据
}

// 新建区块
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block {
	var block = Block{
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Txs:           txs,
	}

	// 生成hash
	// block.setHash()
	// 替换setHash
	// 通过POW生成新的哈希值
	pow := NewProofOfWork(&block)
	// 执行工作量证明算法
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

// 计算区块hash
/*func (b *Block) setHash() {
	// 调用sha256实现hash生成
	// 实现int->hash
	timeStampBytes := Int64ToHex(b.TimeStamp)
	heightBytes := Int64ToHex(b.Height)
	blockBytes := bytes.Join([][]byte{
		timeStampBytes,
		heightBytes,
		b.PrevBlockHash,
		b.Data,
	}, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}
*/
func (b *Block) Serialize() []byte {
	jsonData, err := json.Marshal(&b)
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
		return nil
	}
	return jsonData
}

func Deserialize(jsonData []byte) *Block {
	b := new(Block)
	err := json.Unmarshal(jsonData, &b)
	if err != nil {
		fmt.Println("Error deserializing JSON:", err)

	}
	return b
}

// 把指定区块中所有交易结构都序列化
func (block *Block) HashTransaction() []byte {
	var txsHashes [][]byte
	//将指定区块中所有哈希进行拼接
	for _, tx := range block.Txs {
		txsHashes = append(txsHashes, tx.TxHash)
	}
	txsHash := sha256.Sum256(bytes.Join(txsHashes, []byte{}))
	return txsHash[:]
}
