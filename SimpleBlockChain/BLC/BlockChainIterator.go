package BLC

import (
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
)

type BlockChainIterator struct {
	CurrentHash []byte   // 当前区块hash
	DB          *bolt.DB // 数据库快
}

// 获取区块
func (bcIterator *BlockChainIterator) Next() *Block {
	block := new(Block)
	// 打开数据库并读取
	err := bcIterator.DB.View(func(tx *bolt.Tx) error {
		// 打开数据表
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			// 根据当前的hash获取数据并反序列化
			blockBytes := b.Get(bcIterator.CurrentHash)
			block = DeserializeBlock(blockBytes)
			// 更新当前的hash
			bcIterator.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}
