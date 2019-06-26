package transaction

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestSendRowTxn(t *testing.T) {
	rowTxnBytes, err := CreateRowTxn(common.Client())
	if err != nil {
		t.Fatal(err)
	}
	SendRowTxn(common.Client(), rowTxnBytes)
}
