package whisper

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"log"
)

func Monitor(client *shhclient.Client, keyID string) {
	messages := make(chan *whisperv6.Message)
	criteria := whisperv6.Criteria{
		PrivateKeyID: keyID,
	}
	sub, err := client.SubscribeMessages(context.Background(), criteria, messages)
	if err != nil {
		log.Fatal(err)
	}

	closeCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case message := <-messages:
				fmt.Printf(string(message.Payload)) // "Hello"
				closeCh <- struct{}{}
			}
		}
	}()

	SendMsg(client, keyID)
	<-closeCh
	//	runtime.Goexit() // wait for goroutines to finish
}
