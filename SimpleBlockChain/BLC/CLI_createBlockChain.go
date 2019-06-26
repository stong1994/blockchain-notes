package BLC

func (cli *CLI) createGenesisBlockchain(address, nodeID string) {
	CreateBlockChainWithGenesisBlock(address, nodeID)

	bc := GetBlockchainObject(nodeID)
	defer bc.DB.Close()
	if bc != nil {
		utxoSet := &UTXOSet{bc}
		utxoSet.ResetUTXOSet()
	}
}
