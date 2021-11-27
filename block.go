package main

import (
	"encoding/gob"
	"bytes"
	"log"
	"time"
)

// Block is a single unit in the blockchain
type Block struct {
	Hash  []byte
	Data  []byte
	Prev  []byte
	Nonce int
}


func BuildBlock(data string, Prev []byte) *Block {
	block := &Block{[]byte{}, []byte(data), Prev, 0}
	proof := Proof(block)
	nonce, hash := proof.Compare()

	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

// the first block doesn't have an address to point back to
func Init() *Block {
	return BuildBlock("THE BEGINNING : " + time.Now().Format("2006-01-02 15:04:05"), []byte{})
}

func (b *Block) Encode() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	ErrorHandle(err)

	return res.Bytes()
}

func ErrorHandle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Decode(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	ErrorHandle(err)

	return &block
}
