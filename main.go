package main

import (
	// "bufio"
	"log"
	"fmt"
	"os"
	// "strings"
	"github.com/gin-gonic/gin"

	// server "github.com/hmuir28/go-thepapucoin/p2p"
	"github.com/hmuir28/go-thepapucoin/routes"
	"github.com/hmuir28/go-thepapucoin/database"
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


	// if len(os.Args) != 3 {
	// 	fmt.Println("Usage: go-p2p-server <port> <peer_address>")
	// 	return
	// }

	// port := os.Args[1]               // Port to listen on
	// peerAddress := os.Args[2]        // Address of another peer to connect to
	// peerCh := make(chan server.Peer) // Channel to manage connected peers
	// peers := make([]server.Peer, 0)  // Slice to hold connected peers

	// // Start the server to listen for incoming connections
	// go server.StartServer(port, peerCh)

	// // Connect to an existing peer
	// go server.ConnectToPeer(peerAddress, peerCh)

	// // Start a Goroutine to handle new connections
	// go func() {
	// 	for peer := range peerCh {
	// 		fmt.Println("New peer connected:", peer.Address)
	// 		peers = append(peers, peer)
	// 		go server.HandlePeerConnection(peer, peerCh) // Handle incoming messages
	// 	}
	// }()

	// // Read from stdin to broadcast messages
	// reader := bufio.NewReader(os.Stdin)
	// for {
	// 	fmt.Print("Enter message to broadcast: ")
	// 	message, _ := reader.ReadString('\n')
	// 	message = strings.TrimSpace(message)
	// 	server.BroadcastMessage(peers, message)
	// }

	db := database.NewRedisClient()

	fmt.Println(db)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()

	routes.TransferRoutes(router)

	log.Fatal(router.Run(":" + port))
}
