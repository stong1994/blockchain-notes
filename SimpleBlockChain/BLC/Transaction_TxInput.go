package BLC

import "bytes"

type TxInput struct {
	TxID      []byte // 交易ID
	Vout      int    // 存储Txoutput的vot里面的索引
	Signature []byte // 数字签名
	PublicKey []byte // 公钥 钱包里面
}

// 判断当前的txInput消费,和指定的address是否一致
func (txInput *TxInput) UnLockWithAddress(pubKeyHash []byte) bool {
	//return txInput.ScriptSiq == address
	publicKey := PubKeyHash(txInput.PublicKey)
	return bytes.Compare(pubKeyHash, publicKey) == 0
}
