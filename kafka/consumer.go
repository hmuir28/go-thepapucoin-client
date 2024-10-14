package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Subscriber() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"send-money-topic"}, nil)

	fmt.Println("Waiting for messages...")

	for {
		msg, err := consumer.ReadMessage(-1)
		fmt.Println("pndpowndopwnpodwnpdow")
		fmt.Println(msg)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			continue
		}
		fmt.Printf("Received: %s\n", string(msg.Value))
	}
}

