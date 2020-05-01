package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const subsidy = 100
type Transaction struct {
	ID     []byte
	Vin    []TXInput
	Vout   []TXOutput
}

type TXInput struct {
	TxId   			[]byte
	Vou				int
	ScriptSig		string
}

type TXOutput struct {
	Value   		int
	ScriptPub		string
}

func (in *TXInput) CanUnlockedWith(data string) bool {
	return in.ScriptSig == data
}

func (ou *TXOutput) CanLockedWith(data string) bool {
	return ou.ScriptPub == data
}

func NewCoinBaseTransaction(address, data string) *Transaction {
	in := TXInput{
		TxId:      []byte{},
		Vou:       -1,
		ScriptSig: data,
	}
	ou := TXOutput{
		Value:     subsidy,
		ScriptPub: address,
	}
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{in},
		Vout: []TXOutput{ou},
	}
	tx.SetId()
	return &tx
}

func (tx *Transaction) SetId() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err.Error())
	}
	hash := sha256.Sum256(buff.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].TxId == nil && tx.Vin[0].Vou == -1
}
