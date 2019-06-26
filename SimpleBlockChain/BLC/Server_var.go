package BLC

// 存储节点全局变量
var knowNodes = []string{"localhost:3000"}
var nodeAddress string // 全局变量,节点地址

// 存储哈希
var transactionArray [][]byte
var minerAddress string
var memoryTxPool = make(map[string]*Transaction)
