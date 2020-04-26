package blockchain

import (
	"bitcoin/block"
	"bitcoin/pow"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const bucketName = "block"

type BlockChain struct {
	db            *bolt.DB
	lastHash      []byte
}

type Iterator struct {
	db            *bolt.DB
	cursor 		  []byte
}

func (bc *BlockChain) NewIterator() *Iterator {
	return &Iterator{
		db:     bc.db,
		cursor: bc.lastHash,
	}
}

func (ite *Iterator) Next() *block.Block {
	var blk *block.Block
	err := ite.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		encBlock := b.Get(ite.cursor)
		blk = block.Deserialize(encBlock)
		ite.cursor = blk.PrevHash
		return nil
	})
	if err != nil {
		log.Panic(err.Error())
	}
	return blk
}

func (bc *BlockChain) DB() *bolt.DB {
	return bc.db
}

func NewBlockChain(data string) *BlockChain {
	var lastHash []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err.Error())
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			b, err := tx.CreateBucket([]byte(bucketName))
			if err != nil {
				log.Panic(err.Error())
			}
			gnsBlock := pow.NewBlock(data, nil)
			err = b.Put(gnsBlock.Hash, gnsBlock.Serialize())
			if err != nil {
				log.Panic(err.Error())
			}
			err = b.Put([]byte("lastHash"), gnsBlock.Hash)
			if err != nil {
				log.Panic(err.Error())
			}
			lastHash = gnsBlock.Hash
		} else {
			lastHash = b.Get([]byte("lastHash"))
		}
		return nil
	})
	return &BlockChain{
		db:       db,
		lastHash: lastHash,
	}
}

func (bc *BlockChain) AddBlock(data string) {
	newBlock := pow.NewBlock(data, bc.lastHash)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err.Error())
		}
		err = b.Put([]byte("lastHash"), newBlock.Hash)
		if err != nil {
			log.Panic(err.Error())
		}
		bc.lastHash = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err.Error())
	}
}