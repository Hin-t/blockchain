package BLC

import "crypto/sha256"

// Merkle树实现管理
type MerkleTree struct {
	//根节点
	RootNode *MerkleNode
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
//Merkle根结点之外的其他层必须是偶数个，如果是奇数个，则将最后一个节点复制一份
func NewMerkleTree(txHashes [][]byte) *MerkleTree {
	// 节点列表
	var nodes []*MerkleNode
	//判断交易节点数量，如果是奇数，则拷贝最后一个交易
	if len(txHashes)%2 == 1 {
		txHashes = append(txHashes, txHashes[len(txHashes)-1])
	}
	//遍历所有交易数据，通过哈希生成叶子节点
	for _, data := range txHashes {
		node := MakeMerkleNode(nil, nil, data)
		nodes = append(nodes, node)
	}
	//通过叶子节点创建父节点
	for i := 0; i < len(txHashes)/2; i++ {
		var parentNodes []*MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := MakeMerkleNode(nodes[j], nodes[j+1], nil)
			parentNodes = append(parentNodes, node)
		}
		if len(parentNodes)%2 == 1 {
			parentNodes = append(parentNodes, parentNodes[len(parentNodes)-1])
		}
		//最终只保存了根节点的哈希值
		nodes = parentNodes
	}
	merkleTree := &MerkleTree{nodes[0]}
	return merkleTree
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
	return node
}
