package event_log

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestMonotorLog(t *testing.T) {
	MonotorLog(common.WSClient())
}
