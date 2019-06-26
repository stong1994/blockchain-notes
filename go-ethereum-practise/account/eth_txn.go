package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func ethTxn(client *ethclient.Client, hexKey string, amount *big.Int, gasLimit uint64, gasPrice *big.Int, to string) error {
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	// gas price会根据市场波动,可以用SuggestGasPrice函数,用于根据'x'个先前块来获得平均燃气价格
	gasPrice, err = client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	toAddr := common.HexToAddress(to)

	tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, nil)

	chanID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chanID), privateKey)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	return err
}
