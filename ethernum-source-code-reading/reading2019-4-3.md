## 交易池
交易池之所以存在，一方面为了通过队列的方式缓存交易，防止大量交易的拥堵堵塞；另一方面通过多个队列的方式来筛选不同交易数据，提高交易的处理效率。

在交易池内部实现逻辑中，存在两种类型队列:
- pending队列, 主要用于存储用户创建的准备就绪可执行的交易；
- queue队列，主要用于存储用户创建的暂不可执行的交易

另外，根据以太坊客户端的不同，TxPool数据结构根据客户端的不同分为**全客户端**实现、**轻客户端**两种实现。
两者的区别在于: 轻客户端比全客户端少了queue队列，意味着轻客户端只会按照交易的创建顺序接收本地签名的交易。

**暂不可执行交易**：在以太坊平台中，用户的交易都是按其帐户nonce值进行顺序执行的，任何nonce过大的交易都只能等待前一个(过大nonce值-1)交易完成之后才会将交易从queue队列移至pending状态继续执行。

如果有多个相同nonce值的交易：一般矿工会选择矿工费高的那条交易，而其他交易则会返回失败

## 交易签名


签名流程：
1. 创建交易
创建一条Transcation数据，主要包含九个字段: AccountNonce, Price, GasLimit, Recipient, Amount, Payload, chainId, 0, 0
> 源码位置:core/types/transaction_signing.go:153
```go
// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s EIP155Signer) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.Recipient,
		tx.data.Amount,
		tx.data.Payload,
		s.chainId, uint(0), uint(0),
	})
}
```
2. 序列化数据
根据RLP协议将上述创建的交易进行序列化处理
3. 哈希计算
使用Keccak-256算法计算序列化数据哈希值
4. 私钥签名
使用帐户私钥对上一步哈希数据进行数字签名
5. 字节拆分
对签名结果进行字节拆分，得出R、S 、 V并同步交易实例参数


[EIP155](http://eips.ethereum.org/EIPS/eip-155)鼓励不同的链采用不同的chainID,如下

CHAIN_ID|	Chain(s)
---|---
1|	Ethereum mainnet
2|	Morden (disused), Expanse mainnet
3|	Ropsten
4|	Rinkeby
5|	Goerli
42|	Kovan
1337|	Geth private chains (default)