package main

import (
	"context"

	"github.com/hmuir28/go-thepapucoin/p2p"
	"github.com/hmuir28/go-thepapucoin/database"
)

func main() {

	newInstance := database.NewRedisClient()

	p2pServer := p2p.NewP2PServer()

	var ctx = context.Background()
	
	p2p.StartServer(ctx, p2pServer, newInstance)

	for{}
}
