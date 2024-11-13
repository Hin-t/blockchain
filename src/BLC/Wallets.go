package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

// 钱包集合的管理文件

// 钱包持久化文件
const walletFile = "Wallets.data"

// 实现钱包集合的基本结构
type Wallets struct {
	// key: 地址
	// value: 钱包结构
	Wallets map[string]*Wallet
}

// 初始化钱包集合
func NewWallets() *Wallets {
	// 从钱包文件中获取钱包信息
	// 1. 判断文件是否存在
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets
	}
	// 2. 文件存在，读取内容
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panicf("read the file content failed! %v\n", err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panicf("decode file content  failed! %v\n", err)
	}
	return &wallets
}

// 添加新的钱包到集合中
func (wallets *Wallets) CreateWallet() {
	// 1. 创建钱包
	wallet := NewWallet()
	// 2. 添加
	wallets.Wallets[string(wallet.GetAddress())] = wallet

	// 3. 持久化钱包信息
	wallets.SaveWallet()
}

// 持久化钱包信息(存储到文件中)
func (wallets *Wallets) SaveWallet() {
	var content bytes.Buffer // 钱包内容
	// 注册256椭圆曲线，注册之后可以直接在内部对curve的接口进行编码
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&wallets)
	if err != nil {
		log.Panicf("encode the struct of wallets failed %v !\n", err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panicf("write the content of wallet into file [%s] failed! %v\n", walletFile, err)
	}
}
