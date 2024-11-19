package test

import (
	"blockchain/src/BLC"
	"fmt"
	"testing"
)

func TestNewWallt(t *testing.T) {
	wallet := BLC.NewWallet()
	fmt.Printf("private key : %v \n", wallet.PrivateKey)
	fmt.Printf("public key : %v \n", wallet.PublicKey)
	fmt.Printf("wallet : %v \n", wallet)
}

func TestWallet_GetAddress(t *testing.T) {
	wallet := BLC.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("wallet address: [%s] \n", address)
	fmt.Printf("the validation of current address is %v \n", BLC.IsValidAddress([]byte(address)))
}

func TestBase58Encode(t *testing.T) {

	base58Result := BLC.Base58Encode([]byte("123456789"))
	fmt.Println(base58Result)
}

func TestBase58Decode(t *testing.T) {
	fmt.Println([]byte("123456789"))
	base58Result := BLC.Base58Encode([]byte("123456789"))
	fmt.Println(base58Result)
	fmt.Println(BLC.Base58Decode(base58Result))
}

func TestWallets_CreateWallet(t *testing.T) {
	wallets := BLC.NewWallets()

	wallets.CreateWallet()
	fmt.Printf("wallets: %v \n", wallets.Wallets)

}

// 重置utxo table
func TestResetUTXO(t *testing.T) {
	blockchain := BLC.BlockchainObject()
	defer blockchain.DB.Close()
	utxoset := BLC.UTXOSet{blockchain}
	utxoset.ResetUTXOSet()
}
