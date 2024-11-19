package BLC

import "crypto/sha256"

// Merkle树实现管理
type MerkleTree struct {
	//根节点

}

// Merkle节点结构
type MerkleNode struct {
	//左子节点
	Left *MerkleNode
	//右子节点
	Right *MerkleNode
	//数据
	Data []byte
}

//创建MerkleTree
//txHashes:区块中的哈希交易列表
func NewMerkleTree(txHashes [][]byte) *MerkleTree {

}

//创建Merkle节点
func MakeMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}
	//判断叶子节点
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		//非叶子节点
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}
	//子节点的赋值
	node.Left = left
	node.Right = right
}
