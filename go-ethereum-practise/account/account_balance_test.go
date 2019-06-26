package account

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestGetBalance(t *testing.T) {
	balance, err := getBalance(common.Client(), common.Addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)

}

func TestGetBlockBalance(t *testing.T) {
	balance, err := getBlockBalance(common.Client(), common.Addr, common.BlockNum)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}

func TestGetPendingBalance(t *testing.T) {
	balance, err := getPendingBalance(common.Client(), common.Addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}
