package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

// 区块链文件管理
// 数据库名称
const dbName = "block.db"

// 表名称
const blockTableName = "blocks"

// BlockChain 区块链的基本机构
type BlockChain struct {
	DB  *bolt.DB //区块的切片
	Tip []byte   //保留最新区块的hash值
}

// 判断数据库文件是否存在

// CreateBlockChainWithGenesisBlock 初始化区块链
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {
	if dbExist() {
		fmt.Println("BlockChain already exists")
		os.Exit(1)
	}
	// 生成创世区块
	var blockHash []byte
	// 1. 创建或打开一个数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 创建桶，把生成的创世区块存入数据库中
	err = db.Update(func(tx *bolt.Tx) error {
		// err := tx.DeleteBucket([]byte(blockTableName))

		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panicf("create db [%s] failed %v\n", blockTableName, err)
		}
		//生成一个coinbase交易
		txCoinbase := NewCoinbaseTransaction(address)
		genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
		// 存储
		// 1. key，value分别以什么数据代表--hash
		// 2. 如何把block结构存入到数据库中--序列化
		// fmt.Println("创世区块序列化：", genesisBlock.Serialize())
		err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
		if err != nil {
			log.Panicf("insert the genesisBlock failed %v\n", err)
		}

		blockHash = genesisBlock.Hash
		// 存储最新区块的哈希
		// 1：latest
		err = b.Put([]byte("1"), blockHash)
		if err != nil {
			log.Panicf("save the hash of genesisBlock failed %v\n", err)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return &BlockChain{db, blockHash}
}

func dbExist() bool {
	_, err := os.Stat(dbName)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// AddBlock 添加区块到区块链中
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// newBlock := NewBlock(height, preBlockHash, data)
	// bc.Blocks = append(bc.Blocks, newBlock)

	// 更新区块（insert）
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 1.获取数据库桶
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 2. 获取最新区块哈希值
			lastBlockSequence := b.Get(bc.Tip)
			// 3. 新建区块
			lastBlock := Deserialize(lastBlockSequence)
			newBlock := NewBlock(lastBlock.Height+1, lastBlock.Hash, txs)

			// 4. 存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panicf("insert the genesisBlock failed %v\n", err)
			}
			// 更新最新区块哈希（数据库）
			err = b.Put([]byte("1"), newBlock.Hash)
			if err != nil {
				log.Panicf("updata the hash of newBlock failed %v\n", err)
			}
			// 更新区块链对象中的最新区块哈希
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		return
	}

}

// CreateGenesisBlock 生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

// 遍历数据库，输出所有区块信息
func (bc *BlockChain) PrintChain() {
	fmt.Println("\n打印区块链完整信息...")

	var curBlock *Block
	bcit := bc.Iterator() // 获取迭代器对象

	for {
		fmt.Println("------------------------------------------------------")
		curBlock = bcit.Next()
		fmt.Printf("\tTimeStamp: %v\n", curBlock.TimeStamp)
		fmt.Printf("\tHash: %x\n", curBlock.Hash)
		fmt.Printf("\tPreBlockHash: %x\n", curBlock.PrevBlockHash)
		fmt.Printf("\tHeight: %d\n", curBlock.Height)
		fmt.Printf("\tNonce: %d\n", curBlock.Nonce)
		fmt.Printf("\tTxs: %v\n", curBlock.Txs)
		for _, tx := range curBlock.Txs {
			fmt.Printf("\t\ttx-hash: %x\n", tx.TxHash)
			fmt.Printf("\t\t输入...:\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\t\tvin-hash: %x\n", vin.TxHash)
				fmt.Printf("\t\t\tvin-vout: %v\n", vin.Vout)
				fmt.Printf("\t\t\tvin-PublicKey: %x\n", vin.PublicKey)
				fmt.Printf("\t\t\tvin-Signature: %x\n", vin.Signature)
			}
			fmt.Printf("\t\t输出...: \n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\t\tvout-value: %d\n", vout.Value)
				fmt.Printf("\t\t\tvout-Ripemd160Hash: %x\n", vout.Ripemd160Hash)
			}
		}
		// 退出条件
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PrevBlockHash)
		// 比较
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			// 遍历到创世区块
			break
		}

	}
}

// 获取一个blockchain对象
func BlockchainObject() *BlockChain {
	//获取DB
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("open db [%s] failed %v\n", dbName, err)
	}

	//获取Tip
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	if err != nil {
		log.Panicf("get the blockchain object failed %v\n", dbName, err)
	}
	return &BlockChain{db, tip}
}

// 实现挖矿功能
// 通过接收交易，生成区块
func (blockchain *BlockChain) MineNewBlock(from []string, to []string, amount []string) {
	//搁置交易生成步骤
	var block *Block
	var txs []*Transaction
	// 遍历交易的参与者
	for index, address := range from {
		//调用生成新的交易
		value, _ := strconv.Atoi(amount[index])

		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)

		//追加到txs的交易列表中去
		txs = append(txs, tx)
		//给予交易的发起者（矿工）一定的奖励
		tx = NewCoinbaseTransaction(address)
		txs = append(txs, tx)
	}

	//从数据库中获取最新一个区块
	err := blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//获取最新区块的哈希
			hash := b.Get([]byte("1"))
			//获取最新区块
			blockBytes := b.Get(hash)
			//反序列化
			block = Deserialize(blockBytes)
		}
		return nil
	})
	if err != nil {
		return
	}
	// 此处交易签名验证
	// 对txs中每一笔交易的签名都进行验证
	for _, tx := range txs {
		//只要有一笔交易验证失败。panic
		if blockchain.VerifyTransaction(tx) == false {
			log.Panicf("ERROR: tx [%x] verify failed !%v\n", tx)
		}
	}
	//通过数据库中最近的区块去生成最新的区块（交易打包）
	block = NewBlock(block.Height+1, block.Hash, txs)
	//持久化新生成的区块到数据库中
	err = blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panicf("update the new block to db failed %v\n", err)
			}
			//更新最新区块的哈希值
			err = b.Put([]byte("1"), block.Hash)
			if err != nil {
				log.Panicf("update the latest block hash to db failed %v\n", err)
			}
			blockchain.Tip = block.Hash
		}
		return nil
	})
	if err != nil {
		return
	}
}

// 获取指定地址所有已花费输出
func (blockchain *BlockChain) SpentOutputs(address string) map[string][]int {
	//已花费输出缓存
	spentTXoutputs := make(map[string][]int)
	//获取迭代器对象
	bcit := blockchain.Iterator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			//排除coinbase交易
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.UnLockRipemd160Hash(String2Hash160(address)) {
						key := hex.EncodeToString(in.TxHash)
						//添加到已花费输出的缓存中
						spentTXoutputs[key] = append(spentTXoutputs[key], in.Vout)
					}

				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return spentTXoutputs
}

// 查找指定地址的UTXO
/*
	遍历查找区块链数据库中的每一个区块中的每一个交易
	查找每一个交易的每一个输出
	判断每个输出是否满足下列条件
	1.属于传入的地址
	2.是否未被花费
		1.首先遍历一次数据库，将所有已花费OUTPUT存入一个缓存
		2.再次遍历区块链数据库，检查每一个VOUT是否包含在前面已花费的输出缓存中
*/
func (blockchain *BlockChain) UnUTXOS(address string, txs []*Transaction) []*UTXO {
	//1.遍历数据库，查找所有与address相关的交易
	//获取迭代器
	bcit := blockchain.Iterator()
	var unUTXOS []*UTXO

	// 获取指定地址所有已花费输出
	spentTXOutputs := blockchain.SpentOutputs(address)
	// 缓存迭代
	// 查找缓存中已花费输出
	for _, tx := range txs {
		// 判断coninbaseTransaction
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {

				// 判断用户
				if in.UnLockRipemd160Hash(String2Hash160(address)) {
					// 添加到已花费输出的map中
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}
	// 遍历缓存中的UTXO
	for _, tx := range txs {
		// 添加一个缓存输出的跳转
	WorkCacheTX:
		for index, vout := range tx.Vouts {
			if vout.UnLockScriptPubKeyWithAddress(address) {
				//if vout.CheckPubkeyWithAddress(address) {
				if len(spentTXOutputs) != 0 {
					var isUtxoTx bool // 判断交易是否被其他交易引用
					for txHash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if txHash == txHashStr {
							// 当前遍历到的交易已经有输出被其他的交易的输入所引用
							isUtxoTx = true
							// 添加状态变量，判断指定的output是否被引用
							var isSpentUTXO bool
							for _, voutIndex := range indexArray {
								if index == voutIndex {
									// 该输出被引用
									isSpentUTXO = true
									// 跳出当前vout判断逻辑，进行下一个输出判断
									continue WorkCacheTX
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, vout}
								unUTXOS = append(unUTXOS, utxo)
							}
						}
					}
					if isUtxoTx == false {
						// 说明当前交易中所有与address相关的outputs都是UTXO
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, vout}
					unUTXOS = append(unUTXOS, utxo)
				}
			}
		}
	}
	// 优先遍历缓存中的UTXO，如果余额足够，直接返回，如果不足再遍历db文件中的UTXO

	//数据库迭代，不断获取下一个区块
	for {
		block := bcit.Next()
		//当前地址的未花费输出列表

		//获取指定地址所有已花费输出
		spentTXOutputs := blockchain.SpentOutputs(address)
		//遍历区块中的每笔交易
		for _, tx := range block.Txs {
			//跳转
		work:
			for index, vout := range tx.Vouts {
				//index：当前输出在当前交易的索引位置
				//vout：当前输出
				if vout.UnLockScriptPubKeyWithAddress(address) {
					//if vout.CheckPubkeyWithAddress(address) {
					//当前vout属于传入地址
					if len(spentTXOutputs) != 0 {
						var isSpentOutput bool
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray {
								//txHash：当前输出所引用的交易哈希
								//indexArray:哈希关联的vout索引列表
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									//txHash==tx.TxHash说明当前交易tx至少已经有输出被其他交易的输入引用
									//index ==i 说明正好是当前的输出被其他交易引用
									//跳转到最外层循环，判断下一个VOUT
									isSpentOutput = true
									continue work
								}
							}
						}
						if !isSpentOutput {
							utxo := &UTXO{tx.TxHash, index, vout}
							unUTXOS = append(unUTXOS, utxo)
						}
					} else {
						//将所有输出都添加到未花费输出中
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	fmt.Println("exec the UnUTXOS function")
	return unUTXOS

}

// 查询余额
func (blockchain *BlockChain) getBalance(address string) int {
	var amount int //金额
	utxos := blockchain.UnUTXOS(address, []*Transaction{})
	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}
	return amount
}

// 查找指定地址的可用UTXO,超过amount就中断查找
// 更新当前数据库中指定地址的UTXO数量
// txs:缓存中的交易列表（用于多笔交易处理）
func (blockchain *BlockChain) FindSpendableUTXO(from string, amount int, txs []*Transaction) (int, map[string][]int) {
	// 可用UTXO
	spendableUTXO := make(map[string][]int)

	var value int
	utxos := blockchain.UnUTXOS(from, txs)
	// 遍历UTXO
	for _, utxo := range utxos {
		value += utxo.Output.Value
		// 计算交易哈希
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}
	// 所有的都遍历完成，仍然小于amount
	// 资金不足
	if value < amount {
		fmt.Printf("地址 [%s] 余额不足，当前余额 [%d] \n", from, value)
		os.Exit(1)

	}
	return value, spendableUTXO
}

//通过指定交易hash查找交易
func (blockchain *BlockChain) FindTransaction(ID []byte) Transaction {
	bcit := blockchain.Iterator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(ID, tx.TxHash) == 0 {
				// 找到该交易
				return *tx
			}
		}
		// 退出
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	fmt.Print("没找到交易[%x]\n", ID)
	return Transaction{}
}

// 交易签名
func (blockchain *BlockChain) SignTransaction(tx *Transaction, priKey ecdsa.PrivateKey) {
	//ecdsa
	// coinbase 交易不用签名
	if tx.IsCoinbaseTransaction() {
		return
	}

	//处理交易的input，查找tx中input所引用的vout所属交易（发送者）
	//对所花费的每一笔utxo进行签名
	prevTxs := make(map[string]Transaction)
	for _, vin := range tx.Vins {
		//查找当前交易输入所引用的交易 vin.TxHash
		tx := blockchain.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)] = tx
	}
	//签名
	tx.Sign(priKey, prevTxs)
}

//验证签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}
	prevTxs := make(map[string]Transaction)
	// 查找输入所引用的交易
	for _, vin := range tx.Vins {
		tx := bc.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)] = tx
	}
	//tx.Verify
	return tx.Verify(prevTxs)
}

//查找整条区块链所有已花费输出
func (blockchain *BlockChain) FindAllSpentOutputs() map[string][]*TxINPUT {
	bcit := blockchain.Iterator()
	//存储已花费输出
	spentOutputs := make(map[string][]*TxINPUT)
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			if !tx.IsCoinbaseTransaction() {
				for _, txInput := range tx.Vins {
					txHash := hex.EncodeToString(txInput.TxHash)
					spentOutputs[txHash] = append(spentOutputs[txHash], txInput)
				}
			}
		}
		if isBreakLoop(block.PrevBlockHash) {
			break
		}
	}
	return spentOutputs
}

//退出条件
func isBreakLoop(prevBlockHash []byte) bool {
	var hashInt big.Int
	hashInt.SetBytes(prevBlockHash)
	if hashInt.Cmp(big.NewInt(0)) == 0 {
		return true
	}
	return false
}

// 查找整条区块链中所有地址的UTXO
func (blockchain *BlockChain) FindUTXOMap() map[string]*TXOutputs {
	utxoMaps := make(map[string]*TXOutputs)
	//遍历区块链
	bcit := blockchain.Iterator()
	//查找已花费输出
	spentTXOutputs := blockchain.FindAllSpentOutputs()
	for {
		block := bcit.Next()

		for _, tx := range block.Txs {
			txOutputs := &TXOutputs{[]*TxOutput{}}
			txHash := hex.EncodeToString(tx.TxHash)
			//获取每笔交易的vouts
		WorkOutLoop:
			for index, vout := range tx.Vouts {
				//获取指定交易的输入
				txInputs := spentTXOutputs[txHash]
				if len(txInputs) > 0 {
					isSpent := false
					for _, in := range txInputs {
						//查找指定输出的所有者
						outPubKey := vout.Ripemd160Hash
						inPubKey := in.PublicKey
						if bytes.Compare(Ripemd160Hash(inPubKey), outPubKey) == 0 {
							if index == in.Vout {
								isSpent = true
								continue WorkOutLoop
							}
						}
					}
					if isSpent == false {
						//当前输出没有被包含到txInputs中
						txOutputs.TXOutputs = append(txOutputs.TXOutputs, vout)
					}
				} else {
					//没有input引用该交易的输出，代表则代表当前交易中所有的输出都是UTXO
					txOutputs.TXOutputs = append(txOutputs.TXOutputs, vout)

				}

			}
			utxoMaps[txHash] = txOutputs
		}
		if isBreakLoop(block.PrevBlockHash) {
			break
		}
	}
	return utxoMaps

}
