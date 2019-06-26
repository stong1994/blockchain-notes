package event_log

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestEventReadErc20(t *testing.T) {
	EventReadErc20(common.Client())
}
