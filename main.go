package main

func main() {

	privateKey := loadPrivateKey()

	chain := InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain, &privateKey}

	cli.run()
}
