package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 15

type ProofOfWork struct {
	Block *Block
	Goal  *big.Int
}

// The proof of work to make the blockchain more secure
func Proof(b *Block) *ProofOfWork {
	goal := big.NewInt(1)
	goal.Lsh(goal, uint(256-Difficulty)) //left shift adds zeros to beginning of binary value

	proof := &ProofOfWork{b, goal}

	return proof
}

//ToHex is a utility function that we will use to cast our nonce into a byte
func toHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

//InitData takes your block and adds a nonce (counter/incrementer) to it.
func (proof *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proof.Block.Prev,
			proof.Block.Data,
			toHex(int64(nonce)),
			toHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

//Validate a proof of work to see if it matches the goal
func (proof *ProofOfWork) Validate() bool {
	var hugeInt big.Int

	data := proof.InitData(proof.Block.Nonce)

	hash := sha256.Sum256(data)
	hugeInt.SetBytes(hash[:])

	//this will return true if the hash is valid, and false if not
	return hugeInt.Cmp(proof.Goal) == -1

}

//Compare keeps on testing if the data generated by the proof of work is valid
func (proof *ProofOfWork) Compare() (int, []byte) {
	var hugeInt big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := proof.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hugeInt.SetBytes(hash[:])

		if hugeInt.Cmp(proof.Goal) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()
	return nonce, hash[:]
}
