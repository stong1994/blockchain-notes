package signature

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

func SignatureVerify(rowData string, privateStr string) {
	privateKey, err := crypto.HexToECDSA(privateStr)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	data := []byte(rowData)
	hash := crypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) // 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hexutil.Encode(signature)) // 0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301

	// 检验方法一:恢复公钥,进行比较
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches) // true

	// 检验方法二: SigToPub将返回ECDSA类型的签名公钥,然后比较签名公钥
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println(matches) // true

	// 检验方法三: VerifySignature封装了方法,将公钥切片, 原数据哈希和签名数据传入即可得到结果
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified) // true
}
