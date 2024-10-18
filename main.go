package main

import (
	"github.com/hmuir28/go-thepapucoin/p2p"
	// "github.com/hmuir28/go-thepapucoin/database"
	"github.com/hmuir28/go-thepapucoin/kafka"
)

func main() {
	// Create the blockchain with a genesis block
	// genesisBlock := crypto.GenesisBlock()
	// genesisBlock.Hash = genesisBlock.CalculateHash()
	// blockchain := miner.Blockchain{[]crypto.Block{genesisBlock}}

	// // Simulate some transactions
	// tx1 := crypto.NewTransaction("Alice", "Bob", 10.0)
	// tx2 := crypto.NewTransaction("Bob", "Charlie", 5.0)

	// fmt.Println("Start mining")
	// // Add a new block with proof of work
	// blockchain.AddBlockWithPoW([]crypto.Transaction{tx1, tx2}, 2)
	// fmt.Println("Finished mining")

	// // Print the blockchain
	// for _, block := range blockchain.Blocks {
	// 	fmt.Printf("Index: %d\n", block.Index)
	// 	fmt.Printf("Timestamp: %s\n", block.Timestamp)
	// 	fmt.Printf("PreviousHash: %s\n", block.PreviousHash)
	// 	fmt.Printf("Hash: %s\n", block.Hash)
	// 	fmt.Printf("Nonce: %d\n", block.Nonce)
	// 	fmt.Println("Transactions:")
	// 	for _, tx := range block.Transactions {
	// 		fmt.Printf("\t%s sent %f to %s\n", tx.Sender, tx.Amount, tx.Recipient)
	// 	}
	// 	fmt.Println()
	// }

	p2pServer := p2p.NewP2PServer()

	go kafka.Subscriber(p2pServer)
	
	p2p.StartServer(p2pServer)

	for{}
}
