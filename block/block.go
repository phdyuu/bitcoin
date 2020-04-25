package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	TimeStamp		int64
	Data 			string
	PrevHash		[]byte
	Hash 			[]byte
}

func NewBlock(data string, prevHash []byte) *Block {
	nBlock := &Block{time.Now().Unix(), data, prevHash, nil}
	nBlock.SetHash()
	return nBlock
}

func (b *Block) SetHash() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err.Error())
	}
	hash := sha256.Sum256(buff.Bytes())
	b.Hash = hash[:]
}