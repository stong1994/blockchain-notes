package contract

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestDeployContract(t *testing.T) {
	DeployContract(common.TestClient())
}
