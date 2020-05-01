package main

import (
	"bitcoin/CLI"
	"bitcoin/blockchain"
	"bitcoin/pow"
	"fmt"
	"testing"
)

var address = "phd"
var amount = 30

func TestFunction(t *testing.T) {
	//t.Run("test find utxos", testFindUTXOs)
	t.Run("test pow validate ", testPrintChain)

}

func testFindUTXOs(t *testing.T) {
	bc := blockchain.NewBlockChain(address)
	utxos := bc.FindUTXOs(address)
	accumulated := 0
	for _, out := range utxos {
		accumulated += out.Value
	}
	fmt.Println(accumulated)
}

func testFindSpendableOutputs(t *testing.T) {
	bc := blockchain.NewBlockChain(address)
	spendableTxs, acc := bc.FindSpendableOutputs(address, amount)
	for txId, out := range spendableTxs {
		fmt.Printf("Transactions ID: %s\n", txId)
		fmt.Printf("Out: %v\n", out)
	}
	fmt.Println(acc)
}

func testProofOfWork_Validate(t *testing.T) {
	bc := blockchain.NewBlockChain("phd")
	ite := bc.NewIterator()
	b := ite.Next()
	pw := pow.NewPoW(b)
	fmt.Println(pw.Validate())
}

func testPrintChain(t *testing.T) {
	cli := CLI.CLI{}
	cli.PrintBlockChain()
}