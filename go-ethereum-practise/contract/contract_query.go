package contract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-ethereum-practise/contract/sol"
	"log"
)

func QueryContract(client *ethclient.Client) {
	address := common.HexToAddress("0xd97cf86Cca6429C58C91049d36cF02aE2BF534d3")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}
	version, err := instance.Version(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version) // "1.0"
}
