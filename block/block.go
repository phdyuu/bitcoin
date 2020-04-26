package block

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	TimeStamp		int64
	Data 			string
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