package contract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-ethereum-practise/contract/sol/token"
	"log"
	"math"
	"math/big"
)

// 通过此方法可以获取erc20的数据
func ContractInfo(client *ethclient.Client) {
	tokenAddress := common.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52")
	instance, err := token.NewERC20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0xd183d12ced4accb265b0eda55b3526c7cb102485")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"

	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"
}
