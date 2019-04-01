## cosmos

区块链的架构一般可以分为三层:网络层,共识层和应用层.
在cosmos中,网络层和共识层封装在了tendermint core中,开发者不需要关注内部实现细节,只需要完成自己的应用层即可.

摘抄官方文档的一段话
>  If one wanted to create a Bitcoin-like system on top of ABCI, Tendermint Core would be responsible for
> - Sharing blocks and transactions between nodes  
> - Establishing a canonical/immutable order of transactions (the blockchain)

> The application will be responsible for
> -  Maintaining the UTXO database
> - Validating cryptographic signatures of transactions
> - Preventing transactions from spending non-existent transactions
> - Allowing clients to query the UTXO database.

从这段话中,我们可以看出,**cosmos中的tendermint core就是一个数据库**,我们在应用层所做的东西和我们做传统的后台开发写业务逻辑是没有什么不同的.  
我们在应用层处理好逻辑,然后将对数据库进行CURD操作.只不过这个数据库是去中心化的分布式数据库.

### 共识机制
cosmos的共识机制是基于BFT的**tendermint**,也叫bPOS(bound POS),最多允许**1/3**的故障节点.
![流程图](https://upload-images.jianshu.io/upload_images/1452123-64dc2bd4728a2d9a.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/975/format/webp)

- 流程主要有 NewHeigh -> **Propose** -> **Prevote** -> **Precommit** -> Commit 一共 5 个状态（阶段）。
- 中间三个 Propose -> Prevote -> Precommit是一个round,为共识阶段.
- 在图中可以看到,一个块的提交(commit)可能需要多个round.
- tendermit中锁的机制,保证tendermit不会分叉

#### 推荐文章
- https://www.odaily.com/post/5134145
- https://www.jianshu.com/p/ac82ec874be0
