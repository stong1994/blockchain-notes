package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

func main() {
	client = connectClient()
}

func connectClient() *ethclient.Client {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		panic(err)
	}
	fmt.Println("connect client successful")
	return client
}
