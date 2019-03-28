## 日志
在remix上运行代码，报错
```
transact to SimpleToken.transfer errored: VM error: revert.
revert	The transaction has been reverted to the initial state.
Note: The constructor should be payable if you send value.	Debug the transaction to get more information.
```
remix运行地址：https://remix.ethereum.org/#optimize=true&version=soljson-v0.5.3+commit.10d17f24.js
搜索结果，有人提到有**可能不是代码错误**，在[oraclize](https://dev.oraclize.it/)上运行试试，仍不能执行

测试过程中发现创建合约的地址能够给自己发送代币，但是不能给其它地址发送。回家用ganache试试
