package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

// 用法展示
func PrintUsage() {

	fmt.Println("Usage:")

	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\tgetaccount -- 获取账户列表")
	// 初始化区块链
	fmt.Println("\tcreateblockchain -address address -- 创建区块链")
	// 添加区块
	fmt.Println("\taddblock --data DATA -- 添加区块")
	// 打印完整的区块信息
	fmt.Println("\tprintchain -- 输出区块链信息")
	//通过命令行转账
	fmt.Println("\tsend -from FROM -to TO -amount Amount -- 发起转账")
	fmt.Println("\t转账参数说明")
	fmt.Println("\t\t-from FROM -- 转账源地址")
	fmt.Println("\t\t-to TO -- 转账目标地址")
	fmt.Println("\t\t-amount AMOUBT -- 转账金额")
	fmt.Println("\tgetbalance -address FROM -- 查询指定地址的余额")
	fmt.Println("\t查询余额参数说明")
	fmt.Println("\t\t-address -- 查询余额的地址")

	fmt.Println("\tset_id -port PORT -- 设置端口节点号")
	fmt.Println("\t\t-port -- 访问的节点号")

	fmt.Println("\tstart -- 启动节点服务")

	fmt.Println("\tutxo -test METHOD -- 测试UTXO Table功能中指定的方法")
	fmt.Println("\t\t-METHOD -- 方法名")
	fmt.Println("\t\t\treset -- 重置UTXOtable")
	fmt.Println("\t\t\tbalance -- 查找所有UTXO")
}

// 添加区块
func (cli *CLI) addBlock(txs []*Transaction) {
	//
}

func (cli *CLI) Run() {
	// 获取node id
	nodeID := GetEnvNodeID()
	fmt.Println("Node id is :", nodeID)
	// 检测参数数量
	IsValidArgs()
	// 新建相关命令
	//启动节点命令
	startNodeCmd := flag.NewFlagSet("start", flag.ExitOnError)
	//获取地址命令
	getAccountsCmd := flag.NewFlagSet("getaccounts", flag.ExitOnError)
	// 添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// 创建钱包
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	// 输出区块链完整信息
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 创建区块链
	createBLCWithGenesisiBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	//发起交易
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	//查询余额命令
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	//UTOX测试命令
	utxoTestCmd := flag.NewFlagSet("utxo ", flag.ExitOnError)
	//端口号
	setPortCmd := flag.NewFlagSet("set_id", flag.ExitOnError)

	//-------------------------------------分隔符------------------------------------------

	// 数据参数处理
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")
	//创建区块链时指定矿工奖励地址（接收地址）
	flagCreateBlockChainArg := createBLCWithGenesisiBlockCmd.String("address", "troytan", "指定接收系统奖励")
	//发起交易参数
	flagSendFromArg := sendCmd.String("from", "", "转账源地址")
	flagSendToArg := sendCmd.String("to", "", "转账目标地址")
	flagSendAmountArg := sendCmd.String("amount", "", "转账金额")
	flagGetBalanceArg := getBalanceCmd.String("address", "", "查询余额地址")
	//UTXO测试命令行参数
	flagUTXOTestArg := utxoTestCmd.String("method", "", "UTXO相关操作")
	//端口号参数
	flagSetPortArg := setPortCmd.String("port", "", "端口号")
	// 判断命令
	switch os.Args[1] {
	case "start":
		if err := startNodeCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse cmd of start node failed! %v\n", err)
		}
	case "set_id":
		if err := setPortCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse cmd of set node id failed! %v\n", err)
		}
	case "utxo":
		if err := utxoTestCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse cmd of operate utxo table failed! %v\n", err)
		}
	case "getaccounts":
		if err := getAccountsCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse cmd of  get accounts failed! %v\n", err)
		}
	case "createwallet":
		if err := createWalletCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse cmd of create wallet failed! %v\n", err)
		}
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse getBalanceCmd err: %v", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse sendCmd err: %v", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisiBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse createBLCWithGenesisiBlockCmd err: %v", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printChainCmd err: %v", err)
		}
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse addBlockCmd err: %v", err)
		}
	default:
		// 没有传递任何命令，或者传递的命令不在上面的列表之中
		//fmt.Println("请输入正确的命令")
		PrintUsage()
		os.Exit(1)
	}
	//启动节点服务
	if startNodeCmd.Parsed() {
		cli.startNode(nodeID)
	}
	//设置端口号
	if setPortCmd.Parsed() {
		if *flagSetPortArg == "" {
			fmt.Println("请输入要设置的端口号...")
			//PrintUsage()
			os.Exit(1)
		}
		cli.SetNodeID(*flagSetPortArg)
	}

	// utxo测试命令
	if utxoTestCmd.Parsed() {
		switch *flagUTXOTestArg {
		case "reset":
			cli.TestResetUTXO(nodeID)
		case "balance":
			cli.TestFindUTXOMap()
		}
	}
	// 获取地址列表
	if getAccountsCmd.Parsed() {
		cli.GetAccounts(nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.CreateWallets(nodeID)
	}
	//查询余额
	if getBalanceCmd.Parsed() {
		if *flagGetBalanceArg == "" {
			fmt.Println("请输入查询地址...")
			//PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg, nodeID)
	}
	// 添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		// 调用
		cli.addBlock([]*Transaction{})
	}
	// 输出区块链信息命令
	if printChainCmd.Parsed() {
		cli.printchain(nodeID)
	}
	// 创建区块链命令
	if createBLCWithGenesisiBlockCmd.Parsed() {
		if *flagCreateBlockChainArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockChainArg, nodeID)
	}
	//发起转账
	if sendCmd.Parsed() {
		if *flagSendFromArg == "" {
			fmt.Println("源地址不能为空...")
			PrintUsage()
			os.Exit(1)
		}

		fmt.Printf("\tFROM:[%s] ", JSON2Slice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]", JSON2Slice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]", JSON2Slice(*flagSendAmountArg))

		//fmt.Printf("\tFROM:[%s]\n", *flagSendFromArg)
		//fmt.Printf("\tTO:[%s]\n", *flagSendToArg)
		//fmt.Printf("\tAMOUNT:[%s]\n", *flagSendAmountArg)
		cli.send(JSON2Slice(*flagSendFromArg), JSON2Slice(*flagSendToArg), JSON2Slice(*flagSendAmountArg), nodeID)
	}
}
