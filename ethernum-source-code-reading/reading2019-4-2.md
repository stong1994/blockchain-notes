## 区块与区块头

> 代码位置core/types/block.go:125

```go
// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	header       *Header
	uncles       []*Header
	transactions Transactions

	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}
```

字段|说明
---|---
header| 区块头
uncles|叔块列表，叔块不是不合法区块。而是落后与最长链的分叉区块，系统允许打包六个区块以内的叔块信息并会给予一定奖励
transactions|交易列表。即本身打包确认的转帐或合约交易
hash| TODO
size|当前区块rlp编码后的列表字节数
td|区块总难度值
ReceivedAt| TODO
ReceivedFrom|TODO


> 代码位置core/types/block.go:69

```go
// Header represents a block header in the Ethereum blockchain.
type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        *big.Int       `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`
}
```

字段|字段说明
---|---
ParentHash|	父区块哈希值
UncleHash|	叔块列表哈希值
Coinbase|	挖矿或seal帐户地址
Root|	用户状态树根哈希值
TxHash|	交易列表树根哈希值
ReceiptHash|	交易执行结果树根哈希值
Bloom|	布隆过滤器，快速检索操作日志是否存在
Difficulty|	生产区块难度，取决于具体共识算法
Number|	区块高度，即第多少个区块
GasLimit|	所有交易gaslimit计算总和
GasUsed|	执行所有交易实际消耗的gas，由evm根据实际情况计算而得
Time|	区块创建时间
Extra|	在poa共识中，用于存储sealnode地址
MixDigest|	在POW中由miner进行填充；在POA共识中未做使用
Nonce|	由pow共识运算结果

## 叔块（uncle)
众所周知，在同一时刻可能存在多个矿工同时获得出块权，致使每个矿工在向其他挖矿节点同步区块时出现确认快慢的问题，最终就会存在多个区块并行于当前主链区块，那么这些个并行区块就可以理解为叔块，即Uncle.
![](http://assets.processon.com/chart_image/5ae19d0fe4b039625aef0fed.png)
叔块也是区块，所以结构上和区块一样。

## 交易
以太坊是一个基于账户的系统，目前有两种账户：普通账户和合约账户。每种账户都可以进行交易，因此有两种交易类型：普通交易和合约交易。
- 合约交易：用于发布合约以及调用合约方法。当然，合约交易也是可以向合约帐户进行token转帐
- 普通交易, 用于多个帐户间进行token转帐

> 源码位置:core/types/transaction.go:38
```go
type Transaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}
```
字段|说明
---|---
data|交易数据
hash|交易哈希
size|字节数
from|交易发起方

> 源码位置:core/types/transaction.go:46
```go
type txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}
```
字段|说明
---|---
AccountNonce|账户的nonce值
Price|gasPrice
GasLimit| gasLimit
Recipient|交易接收方，nil代表创建合约
Amount|金额
Payload|TODO
V|TODO
S|TODO
S|TODO
Hash|TODO

如果目标账户是零账户（账户地址是0），交易将创建一个新合约。

创建合约交易的payload（二进制数据）被当作EVM字节码执行。执行的输出做为合约代码被永久存储。这意味着，为了创建一个合约，你不需要向合约发送真正的合约代码，而是发送能够返回真正代码的代码。

### 推荐文章
- 以太坊上发送交易的九种办法：https://www.chainnews.com/articles/817781462334.htm