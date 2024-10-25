package main

import (
	"context"

	"github.com/hmuir28/go-thepapucoin/p2p"
)

func main() {
	p2pServer := p2p.NewP2PServer()

	var ctx = context.Background()
	
	p2p.StartServer(ctx, p2pServer)

	for{}
}
