package block

import (
	"bitcoin/transaction"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	TimeStamp		int64
	Transactions 	[]*transaction.Transaction
	PrevHash		[]byte
	Hash 			[]byte
	Nonce			int64
}

func (b *Block) Serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err.Error())
	}
	return buff.Bytes()
}

func Deserialize(encode []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(encode))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err.Error())
	}
	return &b
}

func (b *Block) HashTransactions() []byte {
	var hashs [][]byte
	for _, tx := range b.Transactions {
		hashs = append(hashs, tx.ID)
	}
	hash := sha256.Sum256(bytes.Join(hashs, []byte{}))
	return hash[:]
}