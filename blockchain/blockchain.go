package blockchain

import (
	"bitcoin/block"
	"bitcoin/pow"
)

type BlockChain struct {
	Blocks 		[]*block.Block
}


func NewBlockChain(data string) *BlockChain {
	gnsBlock := pow.NewBlock(data, nil)
	return &BlockChain{[]*block.Block{gnsBlock}}
}

func (bc *BlockChain) AddBlock(data string) {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	nBlock := pow.NewBlock(data, prevHash)
	bc.Blocks = append(bc.Blocks, nBlock)
}