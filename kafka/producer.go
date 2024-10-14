package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendMessage() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	defer producer.Close()

	topic := "send-money-topic"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Hello Kafka from Golang!"),
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

