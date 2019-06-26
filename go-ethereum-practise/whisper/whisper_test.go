package whisper

import "testing"

func TestSendMsg(t *testing.T) {
	client := WhisperClient()
	keyID := GenKeyPair(client)
	SendMsg(client, keyID)
}

func TestMonitor(t *testing.T) {
	client := WhisperClient()
	keyID := GenKeyPair(client)
	Monitor(client, keyID)
}
