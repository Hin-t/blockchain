package BLC

//当前区块版本信息，是否需要同步
type Version struct {
	//Version     int    //版本号
	Height      int    //当前节点区块高度
	AddressFrom string //当前节点的地址
}
