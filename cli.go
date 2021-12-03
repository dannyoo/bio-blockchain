package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	blockchain *BlockChain
	privateKey *rsa.PrivateKey
}

//printUsage will display what options are availble to the user
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println(" add -fasta <FASTA file location> - add a bio data from .fasta file")
	fmt.Println(" print - prints the blocks in the chain")
	fmt.Println(" wallet - prints your biological datas")
	fmt.Println(" transfer -to <public_key.pem file location> -fasta <FASTA file location> - transfers your biological data")
	fmt.Println(" reset - removes local blockchain database")
}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(0)
	}
}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) reset() {
	if _, err := os.Stat("./db"); !os.IsNotExist(err) {
		err := os.RemoveAll("./db")
		if err != nil {
			log.Fatal(err)
			runtime.Goexit()
		}
		fmt.Println("Succesful Blockchain Database Directory Deletion")
		os.Exit(0)
	}

}

// transfer fasta data with public key
func (cli *CommandLine) transfer(publicKeyFile string, fastaFile string) {
	publicKey := loadPublicKey(publicKeyFile)
	label, seq := readFasta(fastaFile)
	fastaData := label + " " + seq
	encryptedData := rsaEncrypt(fastaData, publicKey)
	cli.blockchain.AddBlock(encryptedData)
	fmt.Println("Transfered fasta data to recipient.")
}

//addBlock allows users to add blocks to the chain via the cli
func (cli *CommandLine) addBlock(file string) {
	label, seq := readFasta(file)
	fastaData := label + " " + seq
	encryptedData := rsaEncrypt(fastaData, cli.privateKey.PublicKey)
	cli.blockchain.AddBlock(encryptedData)
}

// lists data that can be decrypted
func (cli *CommandLine) wallet() {
	counter := 0
	iterator := cli.blockchain.Iteration()
	for {
		block := iterator.Next()
		// attempt to decrypt
		data, err := rsaDecrypt(string(block.Data), *cli.privateKey)
		if err != nil {
			if len(block.Prev) == 0 {
				break
			} else {
				continue
			}
		}
		counter++
		fmt.Println("data: ", data)
		fmt.Printf("block hash: %x\n", block.Hash)
		fmt.Println()
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.Prev) == 0 {
			break
		}
	}
	fmt.Printf("You have %d biological data items on the blockchain.\n", counter)
}

//printChain will display the entire contents of the blockchain
func (cli *CommandLine) printChain() {
	iterator := cli.blockchain.Iteration()

	for {
		block := iterator.Next()
		fmt.Printf("Previous hash: %x\n", block.Prev)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := Proof(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.Prev) == 0 {
			break
		}
	}
}

//run will start up the command line
func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addBlockData := addBlockCmd.String("fasta", "", "FASTA file")
	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	walletCmd := flag.NewFlagSet("wallet", flag.ExitOnError)
	transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	publicKeyFile := transferCmd.String("to", "", "Recipient Public Key File")
	fastaFile := transferCmd.String("fasta", "", "FASTA file data to transfer")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		ErrorHandle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		ErrorHandle(err)

	case "reset":
		err := resetCmd.Parse(os.Args[2:])
		ErrorHandle(err)

	case "wallet":
		err := walletCmd.Parse(os.Args[2:])
		ErrorHandle(err)

	case "transfer":
		err := transferCmd.Parse(os.Args[2:])
		ErrorHandle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}
	// Parsed() will return true if the object it was used on has been called
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if resetCmd.Parsed() {
		cli.reset()
	}
	if walletCmd.Parsed() {
		cli.wallet()
	}
	if transferCmd.Parsed() {
		if *publicKeyFile == "" || *fastaFile == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.transfer(*publicKeyFile, *fastaFile)
	}
}

// Handle erros
func ErrorHandle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
