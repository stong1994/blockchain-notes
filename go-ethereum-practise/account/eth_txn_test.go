package account

import (
	"go-ethereum-practise/common"
	"math/big"
	"testing"
)

func TestEthTxn(t *testing.T) {
	amount := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                 // in units
	gasPrice := big.NewInt(30000000000)       // in wei (30 gwei)
	err := ethTxn(common.Client(), common.HexKey, amount, gasLimit, gasPrice, common.ToAddr)
	if err != nil {
		t.Fatal(err)
	}
}
