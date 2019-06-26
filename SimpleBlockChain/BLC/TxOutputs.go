package BLC

import (
	"bytes"
	"encoding/gob"
	"github.com/labstack/gommon/log"
)

type TxOutputs struct {
	UTXOS []*UTXO
}

func (outs *TxOutputs) Serialize() []byte {
	// 创建buffer
	var result bytes.Buffer
	// 创建编码器
	encoder := gob.NewEncoder(&result)
	// 编码,打包
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeTXOutputs(txOutputBytes []byte) *TxOutputs {
	var txOutputs TxOutputs
	var reader = bytes.NewReader(txOutputBytes)
	// 创建编码器
	decoder := gob.NewDecoder(reader)
	// 解包
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return &txOutputs
}
