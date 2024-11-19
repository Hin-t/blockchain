package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 共识算法管理文件

// 实现POW实例以及相关功能

// 目标难度
const targetBit = 16

// 工作量证明的结构

type ProofOfWork struct {
	// 需要共识验证的区块
	Block *Block
	// 目标难度的哈希
	target *big.Int
}

// 创建一个POW对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// 数据总长度为8位，
	// 需求：需要满足前两位为0，才能解决问题
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{block, target}
}

// 执行POW（），比较哈希值
func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	// 碰撞次数
	var nonce = 0
	var hashInt big.Int
	var hash [32]byte // 生成的哈希值
	// 无限循环，生成符合条件的hash
	for {
		//	生成准备数据
		dataBytes := proofOfWork.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		// 检测生成的hash是否符合条件
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			// 找到了符合条件的hash，中断循环
			break
		}
		nonce++
	}
	fmt.Printf("\n碰撞次数：%d\n", nonce)
	return hash[:], int64(nonce)
}

// 生成准备数据的函数
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	//拼接区块属性，进行hash计算
	// 调用sha256实现hash生成
	// 实现int->hash
	timeStampBytes := Int64ToHex(pow.Block.TimeStamp)
	heightBytes := Int64ToHex(pow.Block.Height)
	data := bytes.Join([][]byte{
		timeStampBytes,
		heightBytes,
		pow.Block.PrevBlockHash,
		pow.Block.HashTransaction(),
		Int64ToHex(nonce),
		Int64ToHex(targetBit),
	}, []byte{})
	hash := sha256.Sum256(data)
	pow.Block.Hash = hash[:]

	return data
}
