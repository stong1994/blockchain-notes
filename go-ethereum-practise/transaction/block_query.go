package transaction

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetLatestBlockNum(client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

func BlockInfoByNum(client *ethclient.Client, num *big.Int) (*types.Block, error) {
	return client.BlockByNumber(context.Background(), num)
}

func TransactionCount(client *ethclient.Client, block *types.Block) (uint, error) {
	return client.TransactionCount(context.Background(), block.Hash())
}
