package event_log

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestReadEventLog(t *testing.T) {
	ReadEventLog(common.WSClient())
}
