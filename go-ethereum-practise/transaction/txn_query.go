package transaction

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func QueryByTx(client *ethclient.Client, tx *types.Transaction) (*types.Message, error) {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func Receipt(client *ethclient.Client, tx *types.Transaction) error {
	// 每个事务都有一个收据，其中包含执行事务的结果，例如任何返回值和日志，以及为“1”（成功）或“0”（失败）的事件结果状态。
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return err
	}
	fmt.Println(receipt.Status, receipt.Logs)
	return nil
}

func QueryByBlockHash(client *ethclient.Client, hash string) {
	blockHash := common.HexToHash(hash)
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tx.Hash().Hex()) // 0xb65985e479047bbe7e06b0799cb6cb6575f1167b1ef8ea52ce0ea53575e97738
	}
}

func QueryByTxHash(client *ethclient.Client, hash string) {
	txHash := common.HexToHash(hash)
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	fmt.Println(isPending)       // false
}
