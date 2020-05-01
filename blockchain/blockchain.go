package blockchain

import (
	"bitcoin/block"
	"bitcoin/pow"
	"bitcoin/transaction"
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const bucketName = "block"
const CoinbaseData = "coinbase"
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


func NewBlockChain(address string) *BlockChain {
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
			coinBaseTx := transaction.NewCoinBaseTransaction(address, CoinbaseData)
			gnsBlock := pow.NewBlock([]*transaction.Transaction{coinBaseTx}, nil)
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

func (bc *BlockChain) AddBlock(transactions []*transaction.Transaction) {
	newBlock := pow.NewBlock(transactions, bc.lastHash)
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

func (bc *BlockChain) FindUnspentTransactions(address string) []*transaction.Transaction {
	var unSpentTXs []*transaction.Transaction
	spentTXs := make(map[string][]int)
	ite := bc.NewIterator()
	for {
		b := ite.Next()
		for _, tx := range b.Transactions {
			txId := hex.EncodeToString(tx.ID)
			Output:
				for outIdx, out := range tx.Vout {
					if spentTXs[txId] != nil {
						for _, spent := range spentTXs[txId] {
							if spent == outIdx {
								continue Output
							}
						}
					}
					if out.CanLockedWith(address) {
						unSpentTXs = append(unSpentTXs, tx)
					}
				}
			if !tx.IsCoinBase() {
				for _, in := range tx.Vin {
					inTxId := hex.EncodeToString(in.TxId)
					if in.CanUnlockedWith(address) {
						spentTXs[inTxId] = append(spentTXs[inTxId], in.Vou)
					}
				}
			}
		}
		if len(b.PrevHash) == 0 {
			break
		}
	}
	return unSpentTXs
}

func (bc *BlockChain) NewUTXOTransaction(from, to string, amount int) *transaction.Transaction {
	unSpentOutputs, acc := bc.FindSpendableOutputs(from, amount)
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput
	if acc < amount {
		log.Panic("Not Enough!")
	}
	for txID, outs := range unSpentOutputs {
		txid, err := hex.DecodeString(txID)
		if err != nil {
			log.Panic(err.Error())
		}
		for _, out := range outs {
			in := transaction.TXInput{
				TxId:      txid,
				Vou:       out,
				ScriptSig: from,
			}
			inputs = append(inputs, in)
		}
	}
	out := transaction.TXOutput{
		Value:     amount,
		ScriptPub: to,
	}
	outputs = append(outputs, out)
	if acc > amount {
		out := transaction.TXOutput{
			Value:     acc - amount,
			ScriptPub: from,
		}
		outputs = append(outputs, out)
	}
	tx := transaction.Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.SetId()
	return &tx
}

func (bc *BlockChain) FindUTXOs(address string) []transaction.TXOutput {
	var unSpentOutputs []transaction.TXOutput
	unSpentTXs := bc.FindUnspentTransactions(address)
	for _, tx := range unSpentTXs {
		for _, out := range tx.Vout {
			if out.CanLockedWith(address) {
				unSpentOutputs = append(unSpentOutputs, out)
			}
		}
	}
	return unSpentOutputs
}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (map[string][]int, int) {
	spendableOutputs := make(map[string][]int)
	unSpentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0
	Spent:
		for _, tx := range unSpentTXs {
			txId := hex.EncodeToString(tx.ID)
			for outIdx, out := range tx.Vout {
				if out.CanLockedWith(address) {
					accumulated += out.Value
					spendableOutputs[txId] = append(spendableOutputs[txId], outIdx)
					if accumulated >= amount {
						break Spent
					}
				}
			}
		}
	return spendableOutputs, accumulated
}