package account

import "testing"

func TestNewWallet(t *testing.T) {
	err := generateWallet()
	if err != nil {
		t.Error(err)
	}
}
