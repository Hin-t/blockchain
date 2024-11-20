package BLC

type Inv struct {
	AddressFrom string   //当前节点地址
	Hashes      [][]byte // 当前节点所拥有的区块哈希列表
}
