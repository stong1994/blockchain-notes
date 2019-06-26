package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) getBalance(address, nodeID string) {
	fmt.Println("查询余额: ", address)
	bc := GetBlockchainObject(nodeID)
	if bc == nil {
		fmt.Println("数据库不存在,无法查询")
		os.Exit(1)
	}
	defer bc.DB.Close()
	//balance := bc.GetBalance(address, []*Transaction{})
	utxoSet := &UTXOSet{bc}
	balance := utxoSet.GetBalance(address)
	fmt.Printf("%s,一共有%d个Token\n", address, balance)
}
