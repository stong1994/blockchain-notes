package BLC

import (
	"bytes"
	"encoding/gob"
	"github.com/labstack/gommon/log"
	"time"
)

// 区块
type Block struct {
	// 字段：
	// 高度Height:其实就是区块的编号，第一个区块叫创世区块，高度为0
	Height int64
	// 上一个区块的哈希值PrevHash
	PrevBlockHash []byte
	// 交易数据Data: 目前先设计[]byte, 后期是Transaction
	//Data []byte
	Txs []*Transaction
	// 时间戳
	TimeStamp int64
	// 哈希值Hash: 32个字节，64个16进制数
	Hash []byte
	// 随机数
	Nonce int64
}

func NewBlock(txs []*Transaction, prevBlockHash []byte, height int64) *Block {
	block := &Block{height, prevBlockHash, txs, time.Now().Unix(), nil, 0}
	pow := NewProofOfWook(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

// 创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, make([]byte, 32, 32), 0)
}

// 将区块序列化,得到一个字节数组
func (block *Block) Serialize() []byte {
	// 创建一个buffer
	var result bytes.Buffer
	// 创建一个编码器
	encoder := gob.NewEncoder(&result)
	// 编码--->打包
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化,得到区块
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	var reader = bytes.NewReader(blockBytes)
	// 创建一个解码器
	decoder := gob.NewDecoder(reader)
	// 解包
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// 将Txs转为byte
func (block *Block) HashTransaction() []byte {
	/*var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]*/
	var txs [][]byte
	for _, tx := range block.Txs {
		txs = append(txs, tx.Serialize())
	}
	mTree := NewMerkleTree(txs)
	return mTree.RootNode.Data
}
