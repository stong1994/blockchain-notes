## 日志

### 外部可调用的函数简述
- name 代币名称
- symbol 代币简称
- totalSupply 代币总量
- owner 创建合约的钱包地址，一般拥有初创合约的所有代币
- INITIAL_SUPPLY 初始的代币总量
- decimals 代币的小数位
- balanceOf 参数为钱包地址，返回余额
> 以上函数都不需要消耗gas

- approve 调用者赋予一个地址 拥有替调用者操作一些代币的权限，比如代币上交易所，就需要给交易所控制一些代币的权限，用于用户提现。需要传参数（使用者和金额）
- decreaseAllowance 减少使用者能够使用的代币的数量。需要传参数（使用者和增加的金额）
- increaseAllowance 增加使用者能够使用的代币的数量。需要传参数（使用者和减少的金额）
- allowance 查看使用者剩余多少某个账户给予的可操作的代币个数。需要传参数（原账户地址和使用者的地址）(不需要消耗gas)
- transfer 交易。需要传参数（接收方和交易的金额）
- transferFrom 使用者调用某个账户的代币进行交易。需要传参（原账户的地址和交易金额）
- burn 调用者燃烧代币，这些代币将永远消失，不但调用者钱包代币减少，代币的总量也会减少。需要参数（燃烧的金额）
- burnFrom 调用者燃烧一部分某个账户给予的可操作的代币。需要参数（原账户地址和燃烧的金额）

### 记录一个错误
在remix上运行合约调用交易的函数，报错
```
transact to SimpleToken.transfer errored: VM error: revert.
revert	The transaction has been reverted to the initial state.
Note: The constructor should be payable if you send value.	Debug the transaction to get more information.
```

原因是remix右上角的控制台中Run界面下的选中的Account，是当前的msg.sender.因为在测试交易的时候需要复制地址，因此会选择合约账户外的钱包地址，此时msg.sender为非构建合约的钱包地址，余额为0，因此交易失败。
