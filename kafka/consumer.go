package kafka

import (
	"context"
	"fmt"
	"log"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/hmuir28/go-thepapucoin/models"
	"github.com/hmuir28/go-thepapucoin/database"
	"github.com/hmuir28/go-thepapucoin/p2p"
)

func Subscriber(p2pServer *p2p.P2PServer) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"send-thepapucoin-topic"}, nil)

	newInstance := database.NewRedisClient()
	var ctx = context.Background()

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			continue
		}

        var transaction models.Transaction
        err = json.Unmarshal(msg.Value, &transaction)
        if err != nil {
            fmt.Printf("Failed to unmarshal message: %v\n", err)
            continue
        }

		database.InsertRecord(ctx, newInstance, transaction)

		// if there are transactions let the peer folks know

		peers := p2pServer.GetPeers()
		
		fmt.Println("--------------- how many peers")
		fmt.Println(peers)
		p2p.BroadcastMessage(peers, "There is a new transaction")
	}
}

