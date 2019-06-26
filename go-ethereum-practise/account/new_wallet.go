package account

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// 生成新钱包
func generateWallet() error {
	// 随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	// 转换为字节
	privateKeyBytes := crypto.FromECDSA(privateKey)
	// 转换为16进制并去掉0x. 得到私钥
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])
	// 得到公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// 去掉0x04
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
	return nil
}
