package transaction

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestErc20Transfer(t *testing.T) {
	Erc20Transfer(common.Client())
}
