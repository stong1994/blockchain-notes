package swarm

import (
	"fmt"
	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
	"log"
)

func SwarmUpload() {
	client := bzzclient.NewClient("http://127.0.0.1:8500")

	file, err := bzzclient.Open("./file/hello.txt")
	if err != nil {
		log.Fatal(err)
	}

	manifestHash, err := client.Upload(file, "", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(manifestHash) //35bbed3ea189e1726775c0553d228ec0c27ebb3e19231a9591a3ed2c5e3fd040
}
