package BLC

func (cli *CLI) TestResetUTXO() {
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	utxoset := UTXOSet{blockchain}
	utxoset.ResetUTXOSet()
}

func (cli *CLI) TestFindUTXOMap() {

}
