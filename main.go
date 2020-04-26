package main

import (
	"bitcoin/blockchain"
	"fmt"
)

func main() {
	bc := blockchain.NewBlockChain("Phd1")
	bc.AddBlock("Phd2")
	bc.AddBlock("Phd3")
	bc.AddBlock("Phd4")
	ite := bc.NewIterator()
	for {
		b := ite.Next()
		fmt.Printf("TimeStamp:    %v\n", b.TimeStamp)
		fmt.Printf("Data:         %s\n", b.Data)
		fmt.Printf("PrevHash:     %x\n", b.PrevHash)
		fmt.Printf("Hash:         %x\n", b.Hash)
		fmt.Printf("Nonce:        %v\n", b.Nonce)

		fmt.Println()
		fmt.Println()
		if len(b.PrevHash) == 0 {
			break
		}
	}
}
