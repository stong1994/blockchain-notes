package account

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	com "go-ethereum-practise/common"
	"regexp"
)

// 如果不符合以太坊正则,返回错误,如果错误为空,并且布尔值是true,为合约地址;如果错误为空,布尔值为false,为用户钱包地址
func AddressValidate(client *ethclient.Client, addr string) (bool, error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	if !re.MatchString(addr) {
		return false, com.NotValidAddress
	}

	// 0x Protocol Token (ZRX) smart contract address
	address := common.HexToAddress(addr)
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	if err != nil {
		return false, err
	}

	isContract := len(bytecode) > 0
	return isContract, nil
}
