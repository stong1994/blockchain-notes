package transaction

import (
	"fmt"
	"go-ethereum-practise/common"
	"math/big"
	"testing"
)

func TestGetLatestBlockNum(t *testing.T) {
	n, err := GetLatestBlockNum(common.Client())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)
}

func TestBlockInfoByNum(t *testing.T) {
	block, err := BlockInfoByNum(common.Client(), big.NewInt(7380501))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time().Uint64())       // 1527211625
	fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println(len(block.Transactions()))   // 144
}

func TestTransactionCount(t *testing.T) {
	block, _ := BlockInfoByNum(common.Client(), big.NewInt(7380501))
	n, err := TransactionCount(common.Client(), block)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)
}
