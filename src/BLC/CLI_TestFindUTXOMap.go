package BLC

func (cli *CLI) TestResetUTXO(nodeID string) {
	blockchain := BlockchainObject(nodeID)
	defer blockchain.DB.Close()
	utxoset := UTXOSet{blockchain}
	utxoset.ResetUTXOSet()
}

func (cli *CLI) TestFindUTXOMap() {

}
