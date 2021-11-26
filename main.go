package main


// TODO: Add comments to codebase...

func main() {

	chain := InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}

    cli.run()
}