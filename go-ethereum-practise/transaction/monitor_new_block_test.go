package transaction

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestMonitorNewBlock(t *testing.T) {
	MonitorNewBlock(common.WSClient())
}
