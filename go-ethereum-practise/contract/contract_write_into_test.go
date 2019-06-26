package contract

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestWriteIntoContract(t *testing.T) {
	WriteIntoContract(common.TestClient())
}
