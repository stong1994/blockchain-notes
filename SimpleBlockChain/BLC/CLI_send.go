package BLC

import (
	"fmt"
	"strconv"
)

func (cli *CLI) send(from, to, amount []string, nodeID string, mineNow bool) {
	blockchain := GetBlockchainObject(nodeID)
	utxoSet := &UTXOSet{blockchain}
	defer blockchain.DB.Close()

	if mineNow {
		blockchain.MineNewBlock(from, to, amount, nodeID)
		// 转账成功,更新
		utxoSet.Update()
	} else {
		// 把交易发送到矿工节点去验证
		fmt.Println("由矿工节点处理...")
		value, _ := strconv.Atoi(amount[0])
		tx := NewSimpleTransaction(from[0], to[0], int64(value), utxoSet, []*Transaction{}, nodeID)

		sendTx(knowNodes[0], tx)
	}
}
