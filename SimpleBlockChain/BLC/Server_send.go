package BLC

import (
	"bytes"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net"
)

func sendVersion(toAddress string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{NODE_VERSION, bestHeight, nodeAddress})

	request := append(commandToBytes(COMMAND_VERSION), payload...)
	sendData(toAddress, request)
}

func sendGetBlocks(toAddress string) {
	payload := gobEncode(GetBlocks{nodeAddress})
	request := append(commandToBytes(COMMAND_GETBLOCKS), payload...)
	fmt.Println("toAddress:%s", toAddress)
	sendData(toAddress, request)
}

// 主节点将自己的所有的区块hash发送给钱包节点
func sendInv(toAddress string, kind string, hashes [][]byte) {
	payload := gobEncode(Inv{nodeAddress, kind, hashes})
	request := append(commandToBytes(COMMAND_INV), payload...)
	sendData(toAddress, request)
}

func sendGetData(toAddress string, kind string, blockHash []byte) {
	payload := gobEncode(GetData{nodeAddress, kind, blockHash})
	request := append(commandToBytes(COMMAND_GETDATA), payload...)
	sendData(toAddress, request)
}

func sendBlock(toAddress string, block []byte) {
	payload := gobEncode(BlockData{nodeAddress, block})
	request := append(commandToBytes(COMMAND_BLOCK), payload...)
	sendData(toAddress, request)
}

func sendTx(toAddress string, tx *Transaction) {
	payload := gobEncode(Tx{nodeAddress, tx})
	request := append(commandToBytes(COMMAND_TX), payload...)
	sendData(toAddress, request)
}

func sendData(to string, data []byte) {
	conn, err := net.Dial(PROTOCOL, to)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()
	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

}
