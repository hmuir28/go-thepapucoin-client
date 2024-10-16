package kafka

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/hmuir28/go-thepapucoin/models"
)

func SendMessage(transaction models.Transaction) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	defer producer.Close()

    transactionBytes, err := json.Marshal(transaction)
    if err != nil {
        log.Fatalf("Error marshaling struct: %v", err)
    }

	topic := "send-thepapucoin-topic"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          transactionBytes,
	}

	err = producer.Produce(message, nil)
	if err != nil {
		log.Fatalf("Failed to send message: %s", err)
	}

	// Wait for delivery report
	e := <-producer.Events()
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		fmt.Printf("Failed to deliver message: %v\n", msg.TopicPartition.Error)
	} else {
		fmt.Printf("Message delivered to %v\n", msg.TopicPartition)
	}
}

