package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func ContractBytecode(client *ethclient.Client) {
	contractAddress := common.HexToAddress("0xd97cf86Cca6429C58C91049d36cF02aE2BF534d3")
	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(bytecode)) // 60806...10029
}
