package contract

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestQueryContract(t *testing.T) {
	QueryContract(common.TestClient())
}
