package account

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"log"
	"os"
)

func createKs(dir, password string) {
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.URL)
	fmt.Println(account.Address.Hex()) // 0x62c37C89dBE8d511c1a8C9EE6fF856e280ED334D
}

func importKs(dir, fileUrl, password string) {
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(fileUrl)
	if err != nil {
		log.Fatal(err)
	}

	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x62c37C89dBE8d511c1a8C9EE6fF856e280ED334D

	if err := os.Remove(fileUrl); err != nil {
		log.Fatal(err)
	}
}
