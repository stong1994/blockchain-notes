package BLC

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
	"os"
)

type UTXOSet struct {
	BlockChain *BlockChain
}

const utxoTableName = "utxoTable"

// 重置数据库表
func (utxoSet *UTXOSet) ResetUTXOSet() {
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic("重置中, 删除表失败")
			}
		}
		b, err := tx.CreateBucket([]byte(utxoTableName))
		if err != nil {
			log.Panic("重置中,创建表失败")
		}
		if b != nil {
			txOutputMap := utxoSet.BlockChain.FindUnSpentOutputMap()
			for txIDStr, outputs := range txOutputMap {
				txID, _ := hex.DecodeString(txIDStr)
				b.Put(txID, outputs.Serialize())
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (utxoSet *UTXOSet) GetBalance(address string) int64 {
	utxos := utxoSet.FindUnSpentOutputsForAddress(address)
	var amount int64
	for _, utxo := range utxos {
		amount += utxo.Output.Value
		fmt.Println(address, "余额", utxo.Output.Value)
	}
	return amount
}

func (utxoSet *UTXOSet) FindUnSpentOutputsForAddress(address string) []*UTXO {
	var utxos []*UTXO
	// 查询数据,遍历所有未消费
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutputs(v)
				for _, utxo := range txOutputs.UTXOS {
					if utxo.Output.UnLockWithAddress(address) {
						utxos = append(utxos, utxo)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return utxos
}

// 用于查询给定地址下的,要转账使用的可以使用的utxo
func (utxoSet *UTXOSet) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	spentableUTXO := make(map[string][]int)
	var total int64 = 0
	// 找出未打包的Transaction中未花费的
	unPackageUTXOS := utxoSet.FindUnPackageSpentableUTXOs(from, txs)
	for _, utxo := range unPackageUTXOS {
		total += utxo.Output.Value
		txIDStr := hex.EncodeToString(utxo.TxID)
		spentableUTXO[txIDStr] = append(spentableUTXO[txIDStr], utxo.Index)
		fmt.Println(amount, ", 未打包, 转账花费:", utxo.Output.Value)
		if total >= amount {
			return total, spentableUTXO
		}
	}
	// 钱不够,找出已经存在数据库中的未花费的
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
		dbLoop:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutputs(v)
				for _, utxo := range txOutputs.UTXOS {
					if utxo.Output.UnLockWithAddress(from) {
						total += utxo.Output.Value
						txIDStr := hex.EncodeToString(utxo.TxID)
						spentableUTXO[txIDStr] = append(spentableUTXO[txIDStr], utxo.Index)
						fmt.Println(amount, ", 数据库,转账花费:", utxo.Output.Value)
						if total > amount {
							break dbLoop
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	if total < amount {
		fmt.Printf("%s, 账户余额不足,不能转账..\n", from)
		os.Exit(1)
	}
	return total, spentableUTXO
}

func (utxoSet *UTXOSet) FindUnPackageSpentableUTXOs(from string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	// 存储已花费
	spentTxOutput := make(map[string][]int)
	for i := len(txs) - 1; i >= 0; i-- {
		unUTXOs = caculate(txs[i], from, spentTxOutput, unUTXOs)
	}
	return unUTXOs
}

// 每次创建区块后,更新未花费的表
func (utxoSet *UTXOSet) Update() {
	/*
		每当创建新区块后,都会花掉一些原来的utxo,产生新的utxo
		删掉已经花费的,增加新产生的未花费
		表中存储的数据结构:
		key: 交易ID
		value: TxInputs
			TxInputs里是UTXO数组
	*/
	// 获取最新的区块,由于该block的产生
	newBlock := utxoSet.BlockChain.Iterator().Next()
	// 遍历该区块的交易
	inputs := []*TxInput{}
	// 未花费
	outsMap := make(map[string]*TxOutputs)
	// 获取已经花费的
	for _, tx := range newBlock.Txs {
		if tx.IsCoinbaseTransaction() {
			continue
		}
		for _, in := range tx.Vins {
			inputs = append(inputs, in)
		}
	}
	fmt.Println("insMap的长度", len(inputs), inputs)
	// 新添加的区块中的未花费的Output
	for _, tx := range newBlock.Txs {
		utxos := []*UTXO{}
	outLoop:
		for index, out := range tx.Vouts {
			isSpent := false
			for _, in := range inputs {
				if bytes.Compare(in.TxID, tx.TxID) == 0 && in.Vout == index && bytes.Compare(out.PubKeyHash, PubKeyHash(in.PublicKey)) == 0 {
					isSpent = true
					continue outLoop
				}
			}
			if isSpent == false {
				utxo := &UTXO{tx.TxID, index, out}
				utxos = append(utxos, utxo)
				fmt.Println("outsMap", out.Value)
			}
		}
		if len(utxos) > 0 {
			txIDStr := hex.EncodeToString(tx.TxID)
			outsMap[txIDStr] = &TxOutputs{utxos}
		}
	}
	fmt.Println("outsMap的长度", len(outsMap), outsMap)

	// 删除已经花费了的
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			// 删除ins中
			for i := 0; i < len(inputs); i++ {
				in := inputs[i]
				fmt.Println(i, "===========================")
				txOutputsBytes := b.Get(in.TxID)
				if len(txOutputsBytes) == 0 {
					continue
				}
				txOutputs := DeserializeTXOutputs(txOutputsBytes)

				// 根据TxID, 如果该txOutputs中已经有output被新区块花掉了,那么将未花掉的添加到utxos里,并标记txoutputs要删除
				// 判断是否需要
				isNeedDelete := false
				utxos := []*UTXO{} // 存储未花费
				for _, utxo := range txOutputs.UTXOS {
					if bytes.Compare(utxo.Output.PubKeyHash, PubKeyHash(in.PublicKey)) == 0 && in.Vout == utxo.Index {
						isNeedDelete = true
					} else {
						utxos = append(utxos, utxo)
					}
				}
				if isNeedDelete == true {
					b.Delete(in.TxID)
					if len(utxos) > 0 {
						txOutputs := &TxOutputs{utxos}
						b.Put(in.TxID, txOutputs.Serialize())
						fmt.Println("删除时:map ", len(outsMap), outsMap)
					}
				}
			}
			// 增加
			for keyID, outPuts := range outsMap {
				keyHashBytes, _ := hex.DecodeString(keyID)
				b.Put(keyHashBytes, outPuts.Serialize())
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
