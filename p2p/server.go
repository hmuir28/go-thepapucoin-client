package p2p

import (
	// "crypto/sha256"
	"context"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"github.com/redis/go-redis/v9"

	"github.com/hmuir28/go-thepapucoin/database"
	"github.com/hmuir28/go-thepapucoin/crypto"
	"github.com/hmuir28/go-thepapucoin/miner"
)

type Peer struct {
	Address string
	Conn    net.Conn
}

type P2PServer struct {
	peers []Peer
}

func (p2pServer P2PServer) GetPeers() []Peer {
	return p2pServer.peers
}

func BroadcastMessage(peers []Peer, message string) {
	for _, peer := range peers {
		_, err := peer.Conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message to peer", peer.Address, ":", err)
			continue
		}
	}
	fmt.Println("Message broadcasted:", message)
}

func NewP2PServer() *P2PServer {
    return &P2PServer{
        peers: []Peer{},
    }
}

func ConnectToPeer(address string, peerCh chan<- Peer) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}

	// Send connection details to peer channel
	peer := Peer{
		Address: address,
		Conn:    conn,
	}
	peerCh <- peer
	fmt.Println("Connected to peer:", address)
}

func HandlePeerConnection(peer Peer, peerCh chan<- Peer) {
	buf := make([]byte, 1024)
	for {
		n, err := peer.Conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closed by", peer.Address)
			peer.Conn.Close()
			peerCh <- peer // Remove peer from the list
			return
		}
		message := string(buf[:n])
		fmt.Printf("Message from %s: %s", peer.Address, message)
	}
}

func SetUpServer(port string, peerCh chan<- Peer) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the incoming connection
		peer := Peer{
			Address: conn.RemoteAddr().String(),
			Conn:    conn,
		}
		peerCh <- peer
	}
}

func StartServer(ctx context.Context, p2pServer *P2PServer, redisClient *redis.Client) {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go-p2p-server <port> <peer_address>")
		return
	}

	port := os.Args[1]               // Port to listen on
	peerAddress := os.Args[2]        // Address of another peer to connect to
	peerCh := make(chan Peer) // Channel to manage connected peers
	// peers := make([]Peer, 0)  // Slice to hold connected peers

	// Start the server to listen for incoming connections
	go SetUpServer(port, peerCh)

	// Connect to an existing peer
	go ConnectToPeer(peerAddress, peerCh)

	// Start a Goroutine to handle new connections
	go func() {
		for peer := range peerCh {
			fmt.Println("New peer connected:", peer.Address)
			p2pServer.peers = append(p2pServer.peers, peer)
			go HandlePeerConnection(peer, peerCh) // Handle incoming messages
		}
	}()

	// Read from stdin to broadcast messages
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to broadcast: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "mine" {
			transactions := database.GetTransactionsInMemory(ctx, redisClient)

			if len(transactions) != 0 {

				// Create the blockchain with a genesis block
				genesisBlock := crypto.GenesisBlock()
				genesisBlock.Hash = genesisBlock.CalculateHash()
				blockchain := miner.Blockchain{[]crypto.Block{genesisBlock}}

				fmt.Println("Start mining")

				fmt.Println(transactions)

				// Add a new block with proof of work
				blockchain.AddBlockWithPoW(transactions, 2)
				fmt.Println("Finished mining")

				// Print the blockchain
				for _, block := range blockchain.Blocks {
					fmt.Printf("Index: %d\n", block.Index)
					fmt.Printf("Timestamp: %s\n", block.Timestamp)
					fmt.Printf("PreviousHash: %s\n", block.PreviousHash)
					fmt.Printf("Hash: %s\n", block.Hash)
					fmt.Printf("Nonce: %d\n", block.Nonce)
					fmt.Println("Transactions:")
					for _, tx := range block.Transactions {
						fmt.Printf("\t%s sent %f to %s\n", tx.Sender, tx.Amount, tx.Recipient)
					}
					fmt.Println()
				}

			}
		} else {
			BroadcastMessage(p2pServer.peers, message)
		}
	}
}
