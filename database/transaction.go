package database

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
    "log"
	"encoding/json"

	"github.com/hmuir28/go-thepapucoin/models"
)

func InsertRecord(ctx context.Context, client *redis.Client, transaction models.Transaction) {

    existingTransactions, err := client.Get(ctx, "transactions").Result()
    var unmarshaledTransactions []models.Transaction

	if err != nil {
		unmarshaledTransactions = []models.Transaction{ transaction }
	} else {
		err = json.Unmarshal([]byte(existingTransactions), &unmarshaledTransactions)

		if err != nil {
			log.Fatalf("Error unmarshaling transactions: %v", err)
		}

		unmarshaledTransactions = append(unmarshaledTransactions, transaction)
	}

	marshaledTransactions, err := json.Marshal(unmarshaledTransactions)

	if err != nil {
        log.Fatalf("Error marshaling transactions: %v", err)
    }

    err = client.Set(ctx, "transactions", marshaledTransactions, 0).Err()

    if err != nil {
        log.Fatalf("Error setting transactions in Redis: %v", err)
    }

    fmt.Println("Transaction stored successfully in Redis!")
}
