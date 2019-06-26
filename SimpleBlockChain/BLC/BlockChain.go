package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
	"math/big"
	"os"
	"strconv"
	"time"
)

// 区块链
type BlockChain struct {
	//Blocks []*Block // 存储有序的区块
	Tip []byte   // 最后区块的hash
	DB  *bolt.DB // 数据库对象
}

// 创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock(address string, nodeID string) {
	DBNAME := fmt.Sprintf(DBNAME, nodeID)
	// 先判断数据库是否存在,如果有,从数据库读取
	if dbExists(DBNAME) {
		fmt.Println("数据库已经存在")
		return
	}
	// 数据库不存在,说明第一次创建,存入数据库中
	fmt.Println("创建创世区块:")
	fmt.Println("数据库不存在...")
	// 创建创世区块
	txCoinBase := NewCoinBaseTransaction(address)
	genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})
	// 打开数据库
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BLOCKTABLENAME))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic("创世区块存储有误...")
			}
			// 存储最新的区块的hash
			b.Put([]byte("1"), genesisBlock.Hash)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//return &BlockChain{genesisBlock.Hash, db}
}

// 添加一个新的区块到区块链中
func (bc *BlockChain) AddBlockToBlockChain(txs []*Transaction) {
	// 创建新区快
	//newBlock := NewBlock(data, prevHash, height)
	// 添加到区块链中
	//bc.Blocks = append(bc.Blocks, newBlock)
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 打开表
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			// 根据最新的hash读取数据,并反序列化最后一个区块
			blockBytes := b.Get(bc.Tip)
			lastBlock := DeserializeBlock(blockBytes)
			// 创建新的区块
			newBlock := NewBlock(txs, lastBlock.Hash, lastBlock.Height+1)
			// 将新的区块序列化并存储
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 更新最后一个哈希值, 以及blockchain的tip
			b.Put([]byte("1"), newBlock.Hash)
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 判断数据库是否存在
func dbExists(DBNAME string) bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}
	return true
}

// 获取一个迭代器的方法
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.Tip, bc.DB}
}

func (bc *BlockChain) PrintChains() {
	/*// 根据bc的tip, 获取最新的hash值,表示当前的hash
	var currentHash = bc.Tip
	// 循环,根据当前的hash读取数据,反序列化得到最后一个区块
	var count = 0
	block := new(Block)
	for {
		err := bc.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BLOCKTABLENAME))
			if b != nil {
				count ++
				fmt.Printf("第%d个区块的信息: \n", count)
				// 获取当前hash对应的数据,并进行反序列化
				blockBytes := b.Get(currentHash)
				block = DeserializeBlock(blockBytes)
				fmt.Printf("\t高度: %d\n", block.Height)
				fmt.Printf("\t上一个区块的hash: %x\n", block.PrevBlockHash)
				fmt.Printf("\t当前的hash: %x\n", block.Hash)

			}
		})
	}*/
	// 获取迭代器对象
	bcIterator := bc.Iterator()
	// 循环迭代
	for {
		block := bcIterator.Next()
		fmt.Printf("第%d个区块的信息: \n", block.Height+1)
		// 获取当前hash对应的数据,并进行反序列化
		fmt.Printf("\t高度: %d\n", block.Height)
		fmt.Printf("\t上一个区块的hash: %x\n", block.PrevBlockHash)
		fmt.Printf("\t当前的hash: %x\n", block.Hash)
		fmt.Println("\t交易: ")
		for _, tx := range block.Txs {
			fmt.Printf("\t\t交易ID: %x\n", tx.TxID)
			fmt.Println("\t\tVins:")
			for _, in := range tx.Vins {
				fmt.Printf("\t\t\tTxID:%x\n", in.TxID)
				fmt.Printf("\t\t\tVout: %d\n", in.Vout)
				fmt.Printf("\t\t\tPublicKey:%v\n", in.PublicKey)
			}
			fmt.Println("\t\tVouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("\t\t\tvalue:%d\n", out.Value)
				fmt.Printf("\t\t\tPubKeyHash:%s\n", out.PubKeyHash)
			}
		}
		fmt.Printf("\t时间: %s \n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("\t次数: %d \n", block.Nonce)

		// 当父hash为0时,说明以遍历到创世区块
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
}

// 获取区块链
func GetBlockchainObject(nodeID string) *BlockChain {
	DBNAME := fmt.Sprintf(DBNAME, nodeID)
	if !dbExists(DBNAME) {
		fmt.Println("数据库不存在,无法读取区块链...")
		return nil
	}
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var blockchain *BlockChain
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			hash := b.Get([]byte("1"))
			blockchain = &BlockChain{hash, db}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return blockchain
}

// 挖掘新的区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string, nodeID string) {
	var txs []*Transaction

	// 奖励
	tx := NewCoinBaseTransaction(from[0])
	txs = append(txs, tx)

	utxoSet := &UTXOSet{bc}
	for i := 0; i < len(from); i++ {
		amountInt, _ := strconv.ParseInt(amount[i], 10, 64)
		tx := NewSimpleTransaction(from[i], to[i], amountInt, utxoSet, txs, nodeID)
		txs = append(txs, tx)
	}

	var block *Block    // 数据库中最后一个block
	var newBlock *Block // 创建的新的Block
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			hash := b.Get([]byte("1"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	// 建立新区块前,对交易进行签名认证
	_txs := []*Transaction{}
	for _, tx := range txs {
		if bc.VerifyTransaction(tx, _txs) != true {
			log.Panic("验证签名失败")
		}
		_txs = append(_txs, tx)
	}

	newBlock = NewBlock(txs, block.Hash, block.Height+1)

	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			b.Put(newBlock.Hash, newBlock.Serialize())
			b.Put([]byte("1"), newBlock.Hash)
			bc.Tip = newBlock.Hash
		}
		return nil
	})
}

// 找到所有未花费的交易输出
func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	/*
		1. 先遍历未打包的交易(参数txs), 找出未花费的Output
		2. 遍历数据库,获取每个块中的Transaction,找出未花费的Output
	*/
	var unUTXOs []*UTXO                      // 未花费
	spentTxOutputs := make(map[string][]int) // 存储已花费
	for i := len(txs) - 1; i >= 0; i-- {
		unUTXOs = caculate(txs[i], address, spentTxOutputs, unUTXOs)
	}

	bcIterator := bc.Iterator()
	for {
		block := bcIterator.Next()
		// 统计未花费
		// 获取block中的每个Transaction
		for i := len(block.Txs) - 1; i >= 0; i-- {
			unUTXOs = caculate(block.Txs[i], address, spentTxOutputs, unUTXOs)
		}
		// 结束迭代
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
	return unUTXOs
}

func (bc *BlockChain) GetBalance(address string, txs []*Transaction) int64 {
	unUTXOs := bc.UnUTXOs(address, txs)
	var amount int64
	for _, utxo := range unUTXOs {
		amount += utxo.Output.Value
	}
	return amount
}

// 转账时查获可用的UTXO
func (bc *BlockChain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	/*
		1.获取所有的UTXO
		2.遍历UTXO
	*/
	var balance int64
	utxos := bc.UnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)
	for _, utxo := range utxos {
		balance += utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxID)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if balance >= amount {
			break
		}
	}
	if balance < amount {
		fmt.Printf("%s 余额不足..总额: %d, 需要: %d \n", from, balance, amount)
		os.Exit(1)
	}
	return balance, spendableUTXO
}

func caculate(tx *Transaction, address string, spentTxOutputs map[string][]int, unUTXOs []*UTXO) []*UTXO {
	// 遍历TxInputs, 表示花费
	if !tx.IsCoinbaseTransaction() {
		for _, in := range tx.Vins {
			// 如果解锁
			fullPayloadHash := Base58Decode([]byte(address))
			pubKeyHash := fullPayloadHash[1 : len(fullPayloadHash)-addressChecksumLen]
			if in.UnLockWithAddress(pubKeyHash) {
				key := hex.EncodeToString(in.TxID)
				spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
			}
		}
	}
	// 遍历TxOutputs
outputs:
	for index, out := range tx.Vouts {
		if out.UnLockWithAddress(address) {
			if len(spentTxOutputs) != 0 {
				var isSpentUTXO bool
				for txID, indexArray := range spentTxOutputs {
					for _, i := range indexArray {
						if i == index && txID == hex.EncodeToString(tx.TxID) {
							isSpentUTXO = true
							continue outputs
						}
					}
				}
				if !isSpentUTXO {
					utxo := &UTXO{tx.TxID, index, out}
					unUTXOs = append(unUTXOs, utxo)
				}
			} else {
				utxo := &UTXO{tx.TxID, index, out}
				unUTXOs = append(unUTXOs, utxo)
			}
		}
	}
	return unUTXOs
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey, txs []*Transaction) {
	if tx.IsCoinbaseTransaction() {
		return
	}
	prevTxs := make(map[string]*Transaction)
	for _, vin := range tx.Vins {
		prevTx := bc.FindTransactionByTxID(vin.TxID, txs)
		prevTxs[hex.EncodeToString(prevTx.TxID)] = prevTx
	}
	tx.Sign(privKey, prevTxs)
}

// 根据交易ID获取对应的Transaction
// 先找未打包的的交易,然后再去每个区块中找
func (bc *BlockChain) FindTransactionByTxID(txID []byte, txs []*Transaction) *Transaction {
	iterator := bc.Iterator()
	for _, tx := range txs {
		if bytes.Compare(txID, tx.TxID) == 0 {
			return tx
		}
	}

	for {
		block := iterator.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(txID, tx.TxID) == 0 {
				return tx
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return &Transaction{}
}

// 验证数字签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction, txs []*Transaction) bool {
	prevTxs := make(map[string]*Transaction)
	for _, vin := range tx.Vins {
		prevTx := bc.FindTransactionByTxID(vin.TxID, txs)
		prevTxs[hex.EncodeToString(prevTx.TxID)] = prevTx
	}
	return tx.Verify(prevTxs)
}

// 查询未花费的Output
func (bc *BlockChain) FindUnSpentOutputMap() map[string]*TxOutputs {
	iterator := bc.Iterator()
	// 存储已经话费
	spentUTXOsMap := make(map[string][]*TxInput)
	// 存储未花费
	unSpentOutputMaps := make(map[string]*TxOutputs)

	for {
		block := iterator.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			txOutputs := &TxOutputs{[]*UTXO{}}
			tx := block.Txs[i]

			if !tx.IsCoinbaseTransaction() {
				for _, txInput := range tx.Vins {
					key := hex.EncodeToString(txInput.TxID)
					spentUTXOsMap[key] = append(spentUTXOsMap[key], txInput)
				}
			}
			txID := hex.EncodeToString(tx.TxID)
		work:
			for index, out := range tx.Vouts {
				txInputs := spentUTXOsMap[txID]
				if len(txInputs) > 0 {
					var isSpent bool
					for _, input := range txInputs {
						inputPubKeyHash := PubKeyHash(input.PublicKey)
						if bytes.Compare(inputPubKeyHash, out.PubKeyHash) == 0 {
							if input.Vout == index {
								isSpent = true
								continue work
							}
						}
					}
					if isSpent == false {
						utxo := &UTXO{tx.TxID, index, out}
						txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxID, index, out}
					txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
				}
			}
			// 设置
			unSpentOutputMaps[txID] = txOutputs
		}
		// 停止迭代
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unSpentOutputMaps
}

// 获取最新区块的高度
func (bc *BlockChain) GetBestHeight() int64 {
	block := bc.Iterator().Next()
	return block.Height
}

// 获取所有的区块的hash
func (bc *BlockChain) GetBlockHashes() [][]byte {
	blockIterator := bc.Iterator()
	var blockHashs [][]byte
	for {
		block := blockIterator.Next()
		blockHashs = append(blockHashs, block.Hash)
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return blockHashs
}

// 根据哈希获取区块
func (bc *BlockChain) GetBlock(blockHash []byte) ([]byte, error) {
	var blockBytes []byte
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			blockBytes = b.Get(blockHash)
		}
		return nil
	})
	return blockBytes, err
}

// 添加区块到数据库
func (bc *BlockChain) AddBlock(block *Block) {
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			blockExist := b.Get(block.Hash)
			if blockExist != nil {
				// 如果存在,不做任何处理
				return nil
			}
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 最新的区块链的hash
			blockHash := b.Get([]byte("1"))
			blockBytes := b.Get(blockHash)
			blockInDB := DeserializeBlock(blockBytes)
			if blockInDB.Height < block.Height {
				b.Put([]byte("1"), block.Hash)
				bc.Tip = block.Hash
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
