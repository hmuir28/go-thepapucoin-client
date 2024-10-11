package p2p

import (
	"fmt"
	"net"
	"os"
)

type Peer struct {
	Address string
	Conn    net.Conn
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

func StartServer(port string, peerCh chan<- Peer) {
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
