package miner

import (
	"github.com/hmuir28/go-thepapucoin/crypto"
)

type Blockchain struct {
	Blocks []crypto.Block
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(transactions []crypto.Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := crypto.CreateBlock(prevBlock, transactions)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) AddBlockWithPoW(transactions []crypto.Transaction, difficulty int) {
	prev := bc.Blocks[len(bc.Blocks)-1]
	new := CreateBlockWithPoW(prev, transactions, difficulty)
	bc.Blocks = append(bc.Blocks, *new.block)
}
