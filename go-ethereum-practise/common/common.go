package common

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Addr     = "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
	BlockNum = 5532993
	HexKey   = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	ToAddr   = "0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d"
)

var (
	NotValidAddress = errors.New("addr is not valid")
)

func Client() *ethclient.Client {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		panic(err)
	}
	fmt.Println("connect client successful")
	return client
}

func WSClient() *ethclient.Client {
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	if err != nil {
		panic(err)
	}
	return client
}

func TestClient() *ethclient.Client {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		panic(err)
	}
	return client
}
