package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Transaction represents a single transaction in the blockchain
type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}

// Block represents each block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PreviousHash string
	Hash         string
	Nonce        int // Add Nonce to keep track of Proof of Work
}

// CalculateHash generates the block's hash
func (b *Block) CalculateHash() string {
	record := string(b.Index) + b.Timestamp + b.PreviousHash + fmt.Sprintf("%d", b.Nonce)
	for _, tx := range b.Transactions {
		record += tx.Sender + tx.Recipient + fmt.Sprintf("%f", tx.Amount)
	}
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// CreateBlock creates a new block
func CreateBlock(prevBlock Block, transactions []Transaction) Block {
	block := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PreviousHash: prevBlock.Hash,
	}
	block.Hash = block.CalculateHash()
	return block
}

// GenesisBlock creates the first block in the blockchain
func GenesisBlock() Block {
	return Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{},
		PreviousHash: "0",
		Hash:         "",
		Nonce:        0,
	}
}

// Create a new transaction
func NewTransaction(sender, recipient string, amount float64) Transaction {
	return Transaction{Sender: sender, Recipient: recipient, Amount: amount}
}
