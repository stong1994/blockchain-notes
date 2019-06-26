package whisper

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"log"
)

func WhisperClient() *shhclient.Client {
	client, err := shhclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func GenKeyPair(client *shhclient.Client) string {
	keyID, err := client.NewKeyPair(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(keyID) // b953d73999d1f9322ab1cf8a078d07731306da276e063725fd1f51a0756cb08c
	return keyID       // 用于加密和解密消息的密钥对
}

func SendMsg(client *shhclient.Client, keyID string) {
	publicKey, err := client.PublicKey(context.Background(), keyID)
	if err != nil {
		log.Print(err)
	}

	fmt.Println(hexutil.Encode(publicKey)) // 0x04133c3a5b92597dc1be733645b0889d904a065ceb78edfc0bf3f1e95d84ce29378245f5e090ec6cdc98385c636cfae37b91ea368cf41c972ed828125190d300df

	message := whisperv6.NewMessage{
		Payload:   []byte("Hello"),
		PublicKey: publicKey, // 加密的公钥
		TTL:       60,        // 消息的活跃时间
		PowTime:   2,         // 做工证明的时间上限
		PowTarget: 2.5,       // 做工证明的时间下限
	}

	// 向网络广播，给它消息，它是否会返回消息的哈希值
	messageHash, err := client.Post(context.Background(), message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messageHash) // 0x229f9a5c4f3a6f9907aff75a789ee0579e393529bece4b63f6dec68d4c2a27e1
}
