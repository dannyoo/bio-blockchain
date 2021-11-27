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
    // fmt.Println(" add -block <BLOCK_DATA> - add a block to the chain")
    fmt.Println(" add -fasta <FASTA file location> - add a bio data from .fasta")
    fmt.Println(" print - prints the blocks in the chain")
    // fmt.Println(" wallet - prints your biological datas")
    // fmt.Println(" transfer -to <public_key> -fasta <FASTA file location> prints your biological datas")
    fmt.Println(" reset - removes local blockchain database")
}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) validateArgs() {
    if len(os.Args) < 2 {
        cli.printUsage()
        //go exit will exit the application by shutting down the goroutine
        // if you were to use os.exit you might corrupt the data
        // runtime.Goexit()
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

//addBlock allows users to add blocks to the chain via the cli
func (cli *CommandLine) addBlock(file string) {
    label, seq := readFasta(file)
    fastaData := label + " " + seq 
    encryptedData := RSA_OAEP_Encrypt(fastaData, cli.privateKey.PublicKey)
    cli.blockchain.AddBlock(string(encryptedData))
    fmt.Println("Added Block!")
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
    resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
    printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
    addBlockData := addBlockCmd.String("fasta", "", "Block data")
    
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
    if resetCmd.Parsed(){
        cli.reset()
    }
}