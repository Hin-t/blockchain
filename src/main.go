package main

import (
	"blockchain/src/BLC"
	"flag"
)

// 启动
// 定义一个字符串变量
var species = flag.String("species", "go", "the usage of flag")

// 定义一个int字符
var num = flag.Int("ins", 1, "ins nums")

func main() {
	// block := BLC.NewBlock(1, nil, []byte("the first block testing"))
	// fmt.Printf("the first block : %v", block)

	//bc := BLC.CreateBlockChainWithGenesisBlock()

	// fmt.Println("blockchain: %s", bc.Blocks[0])
	// 上链
	// AddBlock(height int64, preBlockHash []byte, data []byte)

	//bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1,
	//	bc.Blocks[len(bc.Blocks)-1].Hash,
	//	[]byte("alice send 10 btc to bob"))
	//
	//bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1,
	//	bc.Blocks[len(bc.Blocks)-1].Hash,
	//	[]byte("bob send 10 btc to Alice 1"))

	//for _, block := range bc.Blocks {
	//	fmt.Println("\nblockchain:\n", block)
	//}

	//bc.AddBlock([]byte("a send 100 eth to b"))
	//bc.AddBlock([]byte("b send 100 eth to c"))
	//
	//bc.PrintChain()

	//BLC.SerializeTest()

	// 解析，再flags各种类型参数生效之前，需要对参数进行解析
	//flag.Parse()
	//// 打印参数
	//fmt.Println("a string flag", *species)
	//fmt.Println("ins nums", *num)

	cli := new(BLC.CLI)
	cli.Run()
}
