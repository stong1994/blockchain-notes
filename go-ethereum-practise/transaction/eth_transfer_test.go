package transaction

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestSendTxn(t *testing.T) {
	SendTxn(common.Client())
}
