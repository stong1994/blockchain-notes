package BLC

import "bytes"

// 交易的输出
type TxOutput struct {
	Value int64
	// 一个锁定脚本(ScriptPubKey) 要花这笔钱,必须要解锁该脚本
	PubKeyHash []byte // 公钥
}

// 判断当前的txOutput消费,和指定的address是否一致
func (txOutput *TxOutput) UnLockWithAddress(address string) bool {
	//return txOutput.ScriptPubKey == address
	fullPayloadHash := Base58Decode([]byte(address))
	pubKeyHash := fullPayloadHash[1 : len(fullPayloadHash)-1]
	return bytes.Compare(txOutput.PubKeyHash, pubKeyHash) == 0
}

func NewTxOutput(value int64, address string) *TxOutput {
	txOutput := &TxOutput{value, nil}
	// 设置Ripemd160Hash
	txOutput.Lock(address)
	return txOutput
}

func (txOutput *TxOutput) Lock(address string) {
	publicKeyHash := Base58Decode([]byte(address))
	txOutput.PubKeyHash = publicKeyHash[1 : len(publicKeyHash)-addressChecksumLen]
}
