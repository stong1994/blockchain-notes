package BLC

// 用于表示未花费的
type UTXO struct {
	TxID   []byte // 当前交易ID
	Index  int    // 下表索引
	Output *TxOutput
}
