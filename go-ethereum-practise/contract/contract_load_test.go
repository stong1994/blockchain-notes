package contract

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestLoadContract(t *testing.T) {
	LoadContract(common.TestClient())
}
