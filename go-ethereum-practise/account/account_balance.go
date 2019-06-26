package account

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

// 读取账户余额
func getBalance(client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	return balance, err
}

// 读取账户在某个区块的余额
func getBlockBalance(client *ethclient.Client, addr string, blockNum int64) (*big.Int, error) {
	account := common.HexToAddress(addr)
	blockNumber := big.NewInt(blockNum)
	return client.BalanceAt(context.Background(), account, blockNumber)
}

// 读取待处理的余额
func getPendingBalance(client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	return client.PendingBalanceAt(context.Background(), account)
}
