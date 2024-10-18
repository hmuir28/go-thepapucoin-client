package miner

import (
	"strings"
	"time"

	"github.com/hmuir28/go-thepapucoin/models"
	"github.com/hmuir28/go-thepapucoin/crypto"
)

type Miner struct {
	block *crypto.Block
}

// ProofOfWork adds difficulty to mining new blocks
func (miner *Miner) ProofOfWork(difficulty int) {
	target := strings.Repeat("0", difficulty) // The number of leading zeros required

	for !strings.HasPrefix(miner.block.Hash, target) {
		miner.block.Nonce++ // Increment nonce to change the hash input
		miner.block.Hash = miner.block.CalculateHash()
	}
}

// CreateBlockWithPoW creates a new block with proof of work
func CreateBlockWithPoW(prevBlock crypto.Block, transactions []models.Transaction, difficulty int) Miner {
	block := crypto.Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PreviousHash: prevBlock.Hash,
	}

	miner := Miner{
		block: &block,
	}

	miner.ProofOfWork(difficulty)
	return miner
}
