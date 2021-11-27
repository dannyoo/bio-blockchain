package main

import (
	// you need to run 'go get github.com/dgraph-io/badger'
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./db/chain"
)

//BlockChain is an array of block pointers
type BlockChain struct {
	PrevHash []byte
	Database *badger.DB
}
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// InitBlockChain will be what starts a new blockChain
func InitBlockChain() *BlockChain {
	var prevHash []byte

	options := badger.DefaultOptions(dbPath)
	options.Logger = nil
	db, err := badger.Open(options)
	ErrorHandle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// "previousHash" stand for prev hash
		if _, err := txn.Get([]byte("previousHash")); err == badger.ErrKeyNotFound {
			fmt.Println("Blockchain not found")
			init := Init()
			err = txn.Set(init.Hash, init.Encode())
			ErrorHandle(err)
			err = txn.Set([]byte("previousHash"), init.Hash)
			fmt.Println("Began a New Blockchain")

			prevHash = init.Hash

			return err
		} else {
			item, err := txn.Get([]byte("previousHash"))
			ErrorHandle(err)
			err = item.Value(func(val []byte) error {
				prevHash = val
				return nil
			})
			ErrorHandle(err)
			return err
		}
	})
	ErrorHandle(err)

	blockchain := BlockChain{prevHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("previousHash"))
		ErrorHandle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		ErrorHandle(err)
		return err
	})
	ErrorHandle(err)

	newBlock := BuildBlock(data, lastHash)

	err = chain.Database.Update(func(transaction *badger.Txn) error {
		err := transaction.Set(newBlock.Hash, newBlock.Encode())
		ErrorHandle(err)
		err = transaction.Set([]byte("previousHash"), newBlock.Hash)

		chain.PrevHash = newBlock.Hash
		return err
	})
	ErrorHandle(err)
}

func (chain *BlockChain) Iteration() *BlockChainIterator {
	iteration := BlockChainIterator{chain.PrevHash, chain.Database}

	return &iteration
}

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		ErrorHandle(err)

		err = item.Value(func(val []byte) error {
			block = Decode(val)
			return nil
		})
		ErrorHandle(err)
		return err
	})
	ErrorHandle(err)

	iterator.CurrentHash = block.Prev

	return block
}
