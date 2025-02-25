package p2p

import (
	"context"
	"bufio"
	"fmt"
	"net"
	"os"
	"log"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/hmuir28/go-thepapucoin/models"
	"github.com/hmuir28/go-thepapucoin/crypto"
	"github.com/hmuir28/go-thepapucoin/miner"
)

type Peer struct {
	Address string
	Conn    net.Conn
}

type P2PServer struct {
	Peers []Peer
}

func (p2pServer P2PServer) GetPeers() []Peer {
	return p2pServer.Peers
}

func FindPeerByAddress(peers []Peer, address string) (*Peer, error) {
	for _, peer := range peers {
		fmt.Printf("Consumer Address %s , Peer Address %s \n", address, peer.Address)
		if peer.Address == address {
			return &peer, nil
		}
	}
	return nil, fmt.Errorf("peer with address %s not found", address)
}

func BroadcastMessage(peers []Peer, message string) {
	for _, peer := range peers {
		if peer.Conn == nil {
			continue
		}

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
        Peers: []Peer{},
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

func HandlePeerConnection(ctx context.Context, p2pServer *P2PServer, peer Peer, peerCh chan<- Peer) {
	buf := make([]byte, 1024)
	for {
		if peer.Conn == nil {
			continue
		}

		n, err := peer.Conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closed by", peer.Address)
			peer.Conn.Close()
			peerCh <- peer // Remove peer from the list
			return
		}

		message := string(buf[:n])
		fmt.Printf("Message from %s: %s", peer.Address, message)

		splittedMessage := strings.Split(message, "&")

		trimmedMessage := strings.TrimSpace(splittedMessage[0])

		switch trimmedMessage {
		case "new_transaction":
			var transactions []models.Transaction

			err := json.Unmarshal([]byte(splittedMessage[1]), &transactions)
			if err != nil {
				log.Fatalf("Error decoding transactions: %v", err)
			}

			if len(transactions) != 0 {

				// Create the blockchain with a genesis block
				genesisBlock := crypto.GenesisBlock()
				genesisBlock.Hash = genesisBlock.CalculateHash()
				blockchain := miner.Blockchain{[]crypto.Block{genesisBlock}}

				fmt.Println("Start mining")
 
				strDifficulty := os.Getenv("BLOCK_MINE_DIFFICULTY")

				if strDifficulty == "" {
					strDifficulty = "4"
				}

				difficulty, _ := strconv.Atoi(strDifficulty)

				// Add a new block with proof of work
				blockchain.AddBlockWithPoW(transactions, difficulty)

				fmt.Println("Finished mining")
				
				BroadcastMessage(p2pServer.Peers, "complete_mine")
			}
			break
		default:
			break
		}
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

		fmt.Println(conn.RemoteAddr().String())

		// Handle the incoming connection
		peer := Peer{
			Address: conn.RemoteAddr().String(),
			Conn:    conn,
		}
		peerCh <- peer
	}
}

func StartServer(ctx context.Context, p2pServer *P2PServer) {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go-p2p-server <port> <peer_address>")
		return
	}

	port := os.Args[1]        	// Port to listen on
	peerAddress := os.Args[2]   // Address of another peer to connect to
	peerCh := make(chan Peer) 	// Channel to manage connected peers

	// Start the server to listen for incoming connections
	go SetUpServer(port, peerCh)

	// Connect to an existing peer
	go ConnectToPeer(peerAddress, peerCh)

	// Start a Goroutine to handle new connections
	go func() {
		for peer := range peerCh {
			fmt.Println("New peer connected:", peer.Address)
			p2pServer.Peers = append(p2pServer.Peers, peer)
			go HandlePeerConnection(ctx, p2pServer, peer, peerCh) // Handle incoming messages
		}
	}()

	// Read from stdin to broadcast messages
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to broadcast: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		BroadcastMessage(p2pServer.Peers, message)
	}
}
