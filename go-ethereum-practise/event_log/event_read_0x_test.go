package event_log

import (
	"go-ethereum-practise/common"
	"testing"
)

func TestRead0xEvent(t *testing.T) {
	Read0xEvent(common.Client())
}
