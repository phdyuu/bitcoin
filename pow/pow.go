package pow

import (
	"bitcoin/block"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

const targetBits = 23

type ProofOfWork struct {
	b 			*block.Block
	target		*big.Int
}

func NewPoW(b *block.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b: b, target: target}
}

func (pow *ProofOfWork) PrepareData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		[]byte(strconv.FormatInt(pow.b.TimeStamp, 10)),
		[]byte(pow.b.Data),
		[]byte(strconv.FormatInt(nonce, 10)),
		pow.b.PrevHash,
	}, []byte{})
	return data
}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var nonce int64 = 0
	var hash [32]byte
	hashInt := big.NewInt(1)
	for nonce < math.MaxInt64 {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1{
			fmt.Println("Success")
			fmt.Printf("Hash:   %x\n", hash)
			fmt.Printf("Nonce:  %v\n", nonce)
			fmt.Println()
			break
		}
		nonce++
	}
	return nonce, hash[:]
}

func NewBlock(data string, prevHash []byte) *block.Block {
	b := &block.Block{
		TimeStamp: time.Now().Unix(),
		Data:      data,
		PrevHash:  prevHash,
		Hash:      nil,
		Nonce:     0,
	}
	pw := NewPoW(b)
	nonce, hash := pw.Run()
	b.Nonce = nonce
	b.Hash = hash
	return b
}