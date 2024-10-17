package p2p

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Peer struct {
	Address string
	Conn    net.Conn
}


type Server struct {
	peers []*Peer
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

func StartServer(p2pServer *Server) {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go-p2p-server <port> <peer_address>")
		return
	}

	port := os.Args[1]               // Port to listen on
	peerAddress := os.Args[2]        // Address of another peer to connect to
	peerCh := make(chan Peer) // Channel to manage connected peers
	// peers := make([]Peer, 0)  // Slice to hold connected peers

	// Start the server to listen for incoming connections
	go StartServer(port, peerCh)

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
		BroadcastMessage(p2pServer.peers, message)
	}
}
