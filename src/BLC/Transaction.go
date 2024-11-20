package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
)

//交易管理文件

// 定义一个基本结构
type Transaction struct {
	TxHash []byte      //交易哈希（标识）
	Vins   []*TxINPUT  //输入的列表
	Vouts  []*TxOutput //输出的列表

}

// 实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	//输入
	//coinbase特点
	//txHash：nil
	//vout：
	//ScriptSig
	txInput := &TxINPUT{[]byte{}, -1, nil, nil}
	//输出：
	//value：
	//address:
	//txOutput := &TxOutput{10, String2Hash160(address)}
	txOutput := NewTxOutput(10, address)
	txCoinbase := &Transaction{nil, []*TxINPUT{txInput}, []*TxOutput{txOutput}}
	txCoinbase.HashTransaction()
	return txCoinbase
}

// 生成交易哈希（交易序列化）
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panic("tx Hash failed %V\n", err)
	}
	//添加时间戳标识，不添加会导致所有coinbase交易哈希完全相同
	timestamp := time.Now().UnixNano()
	// 用于生成哈希的原数据
	txHashBytes := bytes.Join([][]byte{result.Bytes(), Int64ToHex(timestamp)}, []byte{})
	//生成哈希值
	hash := sha256.Sum256(txHashBytes)
	tx.TxHash = hash[:]
}

// 生成一个普通转账交易
func NewSimpleTransaction(from string, to string, amount int, bc *BlockChain, txs []*Transaction, nodeID string) *Transaction {
	var txInputs []*TxINPUT
	var txOutputs []*TxOutput

	// 调用可花费UTXO函数
	money, spendableUTXODic := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("money: %v\n", money)
	// 获取钱包集合对象
	wallets := NewWallets(nodeID)
	// 查找对应的钱包结构
	wallet := wallets.Wallets[from]
	//输入
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panicf("decode string to []byte failed %V\n", err)
		}
		// 遍历索引列表
		for _, index := range indexArray {
			txInput := &TxINPUT{txHashBytes, index, wallet.PublicKey, nil}
			txInputs = append(txInputs, txInput)
		}
	}
	//输出(转账源)
	//txOutput := &TxOutput{amount, to}
	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)
	//输出（找零）
	if money > amount {
		//txOutput = &TxOutput{money - amount, from}
		txOutput := NewTxOutput(money-amount, from)
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足...\n")
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction() // 生成一笔完整的交易
	//签名
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

// 判断一个指定的交易是否时一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}

//交易签名
// prevTxs代表当前交易的输入所引用的所有output所属的交易
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	// 处理输入，保证签名的正确性
	// 检查tx每一个输入所引用的交易hash是否包含在prevTxs中
	// 如果没有包含在里面，则说明交易被人篡改过
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("Error : prev transaction is not correct!\n")
		}
	}
	// 提取需要签名的属性
	txCopy := tx.TrimmedCopy()
	// 处理副本的输入
	for vin_id, vin := range txCopy.Vins {
		// 获取关联交易
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 找到发送者（当前输入引用的哈希-输出的哈希）
		txCopy.Vins[vin_id].TxHash = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 生成交易副本哈希
		txCopy.TxHash = txCopy.Hash()
		//调用核心签名函数
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if err != nil {
			log.Panicf("sign to transaction [%x] failed! %v\n", tx, err)
		}

		// 组成交易签名
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[vin_id].Signature = signature
	}

}

// 交易拷贝，生成一个专门用于签名的副本
func (tx *Transaction) TrimmedCopy() Transaction {
	// 从新组装，生成一个新的交易
	var inputs []*TxINPUT
	var outputs []*TxOutput
	// 组装input
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TxINPUT{vin.TxHash, vin.Vout, nil, nil})
	}
	// 组装output
	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TxOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transaction{nil, inputs, outputs}
	return txCopy
}

// 设置用于签名的交易hash
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(tx.Serialize())
	return hash[:]
}

// 交易的序列化
func (tx *Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(tx); err != nil {
		log.Panic("serialize the tx to byte failed! %V\n", err)
	}
	return buffer.Bytes()
}

// 验证签名
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	//检查能否找到相应的交易hash
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("Verify  Error: transaction verify failed!\n")
		}
	}
	// 提取相同的交易签名属性
	txCopy := tx.TrimmedCopy()
	// 使用相同的椭圆
	curve := elliptic.P256()
	//遍历tx输入，对每笔输入引用的输出进行验证
	for vin_id, vin := range tx.Vins {
		// 获取关联交易
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 找到发送者（当前输入引用的哈希-输出的哈希）
		txCopy.Vins[vin_id].TxHash = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 由需要验证的数据生成的交易hash，必须要与签名时使用数据完全一致
		txCopy.TxHash = txCopy.Hash()
		// 在比特币中，签名是一个数值对，r，s代表签名
		// 从输入的signature中获取
		// 获取r，s ，二者长度相等
		r := big.Int{}
		s := big.Int{}
		sigLens := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLens / 2)])
		s.SetBytes(vin.Signature[(sigLens / 2):])
		// 获取公钥
		// 公钥由X，Y组成
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(pubKeyLen / 2)])
		y.SetBytes(vin.PublicKey[(pubKeyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) {
			return false
		}
	}
	// pub *PublicKey, hash []byte, r, s *big.Int
	/*
		// PublicKey represents an ECDSA public key.
		type PublicKey struct {
			elliptic.Curve
			X, Y *big.Int
		}
	*/
	//调用验证签名核心函数
	return true
}
