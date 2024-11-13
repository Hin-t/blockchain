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
}

// 添加区块
func (cli *CLI) addBlock(txs []*Transaction) {
	// 获取到blockchain的对象实例
	if !dbExist() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	blockchain.AddBlock(txs)
}

func (cli *CLI) Run() {
	// 检测参数数量
	IsValidArgs()
	// 新建相关命令
	// 添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// 输出区块链完整信息
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 创建区块链
	createBLCWithGenesisiBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	//发起交易
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	//查询余额命令
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	// 数据参数处理
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")
	//创建区块链时指定矿工奖励地址（接收地址）
	flagCreateBlockChainArg := createBLCWithGenesisiBlockCmd.String("address", "troytan", "指定接收系统奖励")
	//发起交易参数
	flagSendFromArg := sendCmd.String("from", "", "转账源地址")
	flagSendToArg := sendCmd.String("to", "", "转账目标地址")
	flagSendAmountArg := sendCmd.String("amount", "", "转账金额")
	flagGetBalanceArg := getBalanceCmd.String("address", "", "查询余额地址")

	// 判断命令
	switch os.Args[1] {
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

	//查询余额
	if getBalanceCmd.Parsed() {
		if *flagGetBalanceArg == "" {
			fmt.Println("请输入查询地址...")
			//PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg)
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
		cli.printchain()
	}
	// 创建区块链命令
	if createBLCWithGenesisiBlockCmd.Parsed() {
		if *flagCreateBlockChainArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockChainArg)
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
		cli.send(JSON2Slice(*flagSendFromArg), JSON2Slice(*flagSendToArg), JSON2Slice(*flagSendAmountArg))
	}
}
