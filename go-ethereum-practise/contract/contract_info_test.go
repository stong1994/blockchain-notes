package contract

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestContractInfo(t *testing.T) {
	ContractInfo(common.Client())
}
