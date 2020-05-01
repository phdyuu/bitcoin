package CLI

import (
	"bitcoin/blockchain"
	"bitcoin/pow"
	"bitcoin/transaction"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {

}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("    getbalance -address ")
	fmt.Println("    createchain -address")
	fmt.Println("    printchain")
	fmt.Println("    send -from  -to  -amount ")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) getBalance(address string)  {
	var acc int
	bc := blockchain.NewBlockChain(address)
	defer bc.DB().Close()
	utxos := bc.FindUTXOs(address)
	for _, out := range utxos {
		acc += out.Value
	}
	fmt.Printf("balance is : %v\n", acc)
}

func (cli *CLI) createBlockChain(address string) *blockchain.BlockChain {
	return blockchain.NewBlockChain(address)
}

func (cli *CLI) printBlockChain() {
	bc := blockchain.NewBlockChain("")
	defer bc.DB().Close()
	ite := bc.NewIterator()
	for {
		b := ite.Next()
		pw := pow.NewPoW(b)
		fmt.Printf("TimeStamp:    %v\n", b.TimeStamp)
		fmt.Printf("PrevHash:     %x\n", b.PrevHash)
		fmt.Printf("Hash:         %x\n", b.Hash)
		fmt.Printf("Nonce:        %v\n", b.Nonce)
		fmt.Printf("PoW:          %s\n", strconv.FormatBool(pw.Validate()))
		fmt.Println()
		if len(b.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockChain(from)
	defer bc.DB().Close()
	tx := bc.NewUTXOTransaction(from, to, amount)
	bc.AddBlock([]*transaction.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) Run() {
	cli.validateArgs()
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createChainCmd := flag.NewFlagSet("createchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	createChainAddress := createChainCmd.String("address", "", "address")
	getbalanceAddress := getbalanceCmd.String("address", "", "address")
	sendCmdFromAddress := sendCmd.String("from", "", "from")
	sendCmdToAddress := sendCmd.String("to", "", "to")
	sendCmdAmount := sendCmd.Int("amount", 0, "amount")

	switch os.Args[1] {
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
		}
	case "createchain":
		err := createChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if printChainCmd.Parsed() {
		cli.printBlockChain()
	}
	if getbalanceCmd.Parsed() {
		if *getbalanceAddress == "" {
			getbalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceAddress)
	}
	if createChainCmd.Parsed() {
		if *createChainAddress == "" {
			createChainCmd.Usage()
		}
		cli.createBlockChain(*createChainAddress)
	}
	if sendCmd.Parsed() {
		if *sendCmdFromAddress == "" || *sendCmdToAddress == "" || *sendCmdAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendCmdFromAddress, *sendCmdToAddress, *sendCmdAmount)
	}
}