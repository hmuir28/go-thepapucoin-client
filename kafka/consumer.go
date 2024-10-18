package kafka

import (
	"context"
	"fmt"
	"log"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"

	"github.com/hmuir28/go-thepapucoin/models"
	"github.com/hmuir28/go-thepapucoin/database"
	"github.com/hmuir28/go-thepapucoin/p2p"
)

func Subscriber(ctx context.Context, p2pServer *p2p.P2PServer, redisClient *redis.Client) {
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

		database.InsertRecord(ctx, redisClient, transaction)

		peers := p2pServer.GetPeers()
		
		p2p.BroadcastMessage(peers, "There is a new transaction")
	}
}

