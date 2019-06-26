package swarm

import (
	"fmt"
	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
	"io/ioutil"
	"log"
)

func SwarmDownload() {
	client := bzzclient.NewClient("http://127.0.0.1:8500")
	manifestHash := "35bbed3ea189e1726775c0553d228ec0c27ebb3e19231a9591a3ed2c5e3fd040"
	manifest, isEncrypted, err := client.DownloadManifest(manifestHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(isEncrypted) // false

	for _, entry := range manifest.Entries {
		fmt.Println(entry.Hash)        // 92672a471f4419b255d7cb0cf313474a6f5856fb347c5ece85fb706d644b630f
		fmt.Println(entry.ContentType) // text/plain; charset=utf-8
		fmt.Println(entry.Size)        // 11
		fmt.Println(entry.Path)        // ""
	}

	file, err := client.Download(manifestHash, "")
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content)) // hello world
}
