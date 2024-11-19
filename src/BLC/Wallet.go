package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// 钱包管理相关文件

// 校验和长度
const addressChecksumLen = 4

type Wallet struct {
	// 1. 私钥
	PrivateKey ecdsa.PrivateKey
	// 2. 公钥
	PublicKey []byte
}

// 创建一个钱包
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{PrivateKey: privateKey, PublicKey: publicKey}
}

// 通过钱包生成公钥-私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 1.获取一个椭圆
	curve := elliptic.P256()
	// 2.通过椭圆相关算法生成私钥
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic("ecdsa generate private key failed! %v\n", err)
	}
	// 3. 通过私钥生成公钥
	pubKey := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)

	return *priv, pubKey
}

// 生成地址

// 实现双哈希
func Ripemd160Hash(pubKey []byte) []byte {
	// 1. sha256
	hash256 := sha256.New()
	hash256.Write(pubKey)
	hash := hash256.Sum(nil)
	// 2. rimpemd160
	rmd160 := ripemd160.New()
	rmd160.Write(hash)
	return rmd160.Sum(nil)

}

// 生成校验和
func CheckSum(input []byte) []byte {
	firstHash := sha256.Sum256(input)
	seondHash := sha256.Sum256(firstHash[:])
	return seondHash[:addressChecksumLen]
}

// 通过钱包（公钥）获取地址
func (w *Wallet) GetAddress() []byte {
	// 1. 获取hash160
	ripmed160Hash := Ripemd160Hash(w.PublicKey)
	// 2. 获取校验和
	checkSumBytes := CheckSum(ripmed160Hash)
	// 3. 地址组成成员拼接
	addressBytes := append(ripmed160Hash, checkSumBytes...)
	// 4. base58编码
	b58Bytes := Base58Encode(addressBytes)

	return b58Bytes
}

// 判断地址有效性
func IsValidAddress(addressBytes []byte) bool {
	// 1. 地址通过Base58Decode进行解码
	public_CheckSumByte := Base58Decode(addressBytes)
	// 2. 拆分，进行校验和校验
	checkSumBytes := public_CheckSumByte[len(public_CheckSumByte)-addressChecksumLen:]
	// 3. 生成
	rimpemd160hash := public_CheckSumByte[:len(public_CheckSumByte)-addressChecksumLen]
	checkBytes := CheckSum(rimpemd160hash)
	// 4. 校验
	if bytes.Equal(checkBytes, checkSumBytes) {
		return true
	}
	return false
}
