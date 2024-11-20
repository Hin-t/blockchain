package BLC

// 网络服务常量管理

// 协议
const PROTOCOL = "tcp"

// 命令长度
const COMMAND_LENGTH = 12

// 命令分类
const (
	// 验证当前节点末端是否是最新区块
	CMD_VERSION = "version"
	// 从最长链获取区块
	CMD_GETBLOCKS = "getblocks"
	// 向其他节点展示当前节点有哪些区块
	CMD_INV = "inv"
	// 请求指定区块
	CMD_GETDATA = "getdata"
	// 接收到新区块进行处理
	CMD_BLOCK = "block"
)
