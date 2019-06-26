package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const TargetBit = 16 // 256位的Hash里面前面至少有16个零

type ProofOfWork struct {
	// 要验证的区块
	Block *Block
	// 大整数存储,目标哈希
	Target *big.Int
}

func NewProofOfWook(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// 左移256个bit位
	target = target.Lsh(target, 256-TargetBit)
	return &ProofOfWork{block, target}
}

// 返回有效的哈希和nonce值
func (pow *ProofOfWork) Run() ([]byte, int64) {
	nonce := 0
	hashInt := new(big.Int)
	var hash [32]byte
	for {
		// 获取字节数组
		dataBytes := pow.prepareData(nonce)
		// 生成hash
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//fmt.Println(hash)
		fmt.Printf("\r%d:%x", nonce, hash)
		if pow.Target.Cmp(hashInt) == 1 {
			break
		}
		nonce++
	}
	fmt.Println()
	return hash[:], int64(nonce)
}

// 根据block生成一个byte数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransaction(),
			IntToHex(pow.Block.TimeStamp),
			IntToHex(int64(TargetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) IsValid() bool {
	hashInt := new(big.Int)
	hashInt.SetBytes(pow.Block.Hash)
	return pow.Target.Cmp(hashInt) == 1
}
