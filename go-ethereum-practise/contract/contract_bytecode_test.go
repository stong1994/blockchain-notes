package contract

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestContractBytecode(t *testing.T) {
	ContractBytecode(common.TestClient())
}
