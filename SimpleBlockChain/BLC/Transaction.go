package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"math/big"
	"time"
)

type Transaction struct {
	TxID  []byte     // 交易ID
	Vins  []*TxInput // 输入
	Vouts []*TxOutput
}

/*
Transaction 分两种情况
1.创世区块创建时的Transaction
2.转账时产生的Transaction
*/
func NewCoinBaseTransaction(address string) *Transaction {
	txInput := &TxInput{[]byte{}, -1, nil, []byte{}}
	txOutput := NewTxOutput(10, address)
	txCoinbase := &Transaction{[]byte{}, []*TxInput{txInput}, []*TxOutput{txOutput}}
	txCoinbase.SetTxID()
	return txCoinbase
}

// 设置交易ID 也就是hash
func (tx *Transaction) SetTxID() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	bufferBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()), buff.Bytes()}, []byte{})
	hash := sha256.Sum256(bufferBytes)
	tx.TxID = hash[:]
}

func NewSimpleTransaction(from, to string, amount int64, utxoSet *UTXOSet, txs []*Transaction, nodeID string) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	//balance, spendableUTXO := bc.FindSpendableUTXOs(from, amount, txs)
	balance, spendableUTXO := utxoSet.FindSpendableUTXOs(from, amount, txs)
	// 获取钱包
	wallets := NewWallets(nodeID)
	wallet := wallets.WalletsMap[from]
	// 代表消费
	for txID, indexArray := range spendableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TxInput{txIDBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}
	// 转账
	txOutput1 := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput1)
	// 找零
	txOutput2 := NewTxOutput(balance-amount, from)
	txOutputs = append(txOutputs, txOutput2)
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	// 设置hash
	tx.SetTxID()
	// 进行签名
	utxoSet.BlockChain.SignTransaction(tx, wallet.PrivateKey, txs)
	return tx
}

// 判断当前交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}

// 签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]*Transaction) {
	// coinbase交易,无需签名
	if tx.IsCoinbaseTransaction() {
		return
	}

	// input没有对应的transaction, 无法签名
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString((vin.TxID))].TxID == nil {
			log.Panic("当前的input没有对应的Transaction")
		}
	}
	// 获取Transaction的部分数据的副本
	txCopy := tx.TrimmedCopy()
	for index, input := range txCopy.Vins {
		prevTx := prevTxs[hex.EncodeToString(input.TxID)]
		// 为txCopy设置新的交易ID, txID -> []byte{}, Vount, sign -->nil, publicKey --> 对应输出的公钥哈希
		input.Signature = nil                                 // 双保险
		input.PublicKey = prevTx.Vouts[input.Vout].PubKeyHash // 设置input的公钥的对应输出的公钥哈希
		data := txCopy.getData()                              // 设置新的txID
		input.PublicKey = nil                                 // 将publicKey 置为nil

		// 签名
		// 通过privKey 对 txCopy.ID 进行签名.一个ECDSA签名就是一对数字,将数字连接起来,并存储到输出的Signature字段
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, data)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[index].Signature = signature
	}
}

// 获取签名所需要的Transaction的副本
// 创建tx的副本,需要剪裁数据
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput
	for _, input := range tx.Vins {
		inputs = append(inputs, &TxInput{input.TxID, input.Vout, nil, nil})
	}
	for _, output := range tx.Vouts {
		outputs = append(outputs, &TxOutput{output.Value, output.PubKeyHash})
	}
	txCopy := Transaction{tx.TxID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) Serialize() []byte {
	jsonByte, err := json.Marshal(tx)
	if err != nil {
		log.Panic(err)
	}
	return jsonByte
}

func (tx Transaction) getData() []byte {
	txCopy := tx
	txCopy.TxID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// 验证数字签名
func (tx *Transaction) Verify(prevTxs map[string]*Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}
	// input没有对应的Transaction,无法签名
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxID)].TxID == nil {
			log.Panic("当前的input没有对应的Transaction,无法验证")
		}
	}
	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()
	for index, input := range tx.Vins {
		prevTx := prevTxs[hex.EncodeToString(input.TxID)]
		txCopy.Vins[index].Signature = nil
		txCopy.Vins[index].PublicKey = prevTx.Vouts[input.Vout].PubKeyHash
		data := txCopy.getData()
		txCopy.Vins[index].PublicKey = nil
		// 签名中的s和r
		r := big.Int{}
		s := big.Int{}
		sigLen := len(input.Signature)
		r.SetBytes(input.Signature[:sigLen/2])
		s.SetBytes(input.Signature[sigLen/2:])

		// 通过公钥,产生新的s和r,与原来的进行对比
		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PublicKey)
		x.SetBytes(input.PublicKey[:keyLen/2])
		y.SetBytes(input.PublicKey[keyLen/2:])

		// 根据椭圆曲线,以及x,y 获取公钥
		// 我们使用从输入提取的公钥新建了一个ecdsa.PublicKey
		rawPublicKey := ecdsa.PublicKey{curve, &x, &y}
		// 一个签名就是一对数字, 一个公钥就是一对坐标
		if ecdsa.Verify(&rawPublicKey, data, &r, &s) == false {
			return false
		}
	}
	return false
}
