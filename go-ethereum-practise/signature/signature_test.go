package signature

import "testing"

func TestSignatureGenerate(t *testing.T) {
	SignatureGenerate("hello", "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
}

func TestSignatureVerify(t *testing.T) {
	SignatureVerify("hello", "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
}
