package main

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Hash  []byte
	Data  []byte
	Prev  []byte
	Nonce int
}

// creates a standard block
func BuildBlock(data string, Prev []byte) *Block {
	block := &Block{[]byte{}, []byte(data), Prev, 0}
	proof := Proof(block)
	nonce, hash := proof.Compare()

	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

// Create the initial block
func Init() *Block {
	return BuildBlock("THE BEGINNING : "+time.Now().Format("2006-01-02 15:04:05"), []byte{})
}

// Encodes the data to stored on badger db
func (b *Block) Encode() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	ErrorHandle(err)

	return res.Bytes()
}

// Decodes the data to stored on badger db into a struct
func Decode(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	ErrorHandle(err)

	return &block
}