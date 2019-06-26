package transaction

import (
	"go-ethereum-practise/common"
	"math/big"
	"testing"
)

func TestQueryByTx(t *testing.T) {
	block, err := BlockInfoByNum(common.Client(), big.NewInt(7380501))
	if err != nil {
		t.Fatal(err)
	}
	if block.Transactions().Len() <= 0 {
		t.Fatal("invalid block")
	}
	// 此tx只能获取to地址,不能获取from地址
	tx := block.Transaction(block.Transactions()[0].Hash())
	t.Log(tx.To().Hex())

	msg, err := QueryByTx(common.Client(), tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msg.From().Hex(), msg.To().Hex())
}

func TestReceipt(t *testing.T) {
	block, err := BlockInfoByNum(common.Client(), big.NewInt(7380501))
	if err != nil {
		t.Fatal(err)
	}
	if block.Transactions().Len() <= 0 {
		t.Fatal("invalid block")
	}
	tx := block.Transaction(block.Transactions()[0].Hash())
	err = Receipt(common.Client(), tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryByBlockHash(t *testing.T) {
	QueryByBlockHash(common.Client(), "0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
}

func TestQueryByTxHash(t *testing.T) {
	QueryByTxHash(common.Client(), "0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2")
}
