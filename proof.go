package main

import (
	"crypto/sha256"
	"log"
	"fmt"
	"math"
	"math/big"
	"encoding/binary"
	"bytes"
)


const Difficulty = 12

type ProofOfWork struct {
	Block *Block
	Goal  *big.Int
}

func Proof(b *Block) *ProofOfWork {
	goal := big.NewInt(1)
	goal.Lsh(goal, uint(256-Difficulty)) //left shift

	proof := &ProofOfWork{b, goal}

	return proof
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

//Validate will check our Run() function performed as expected
func (proof *ProofOfWork) Validate() bool {
	var hugeInt big.Int

	data := proof.InitData(proof.Block.Nonce)

	hash := sha256.Sum256(data)
	hugeInt.SetBytes(hash[:])

	//this will return true if the hash is valid, and false if not
	return hugeInt.Cmp(proof.Goal) == -1

}

//Compare will hash our data, turn that hash into a big int, and then compare that big int to our Goal which is inside  the Proof Of Work Struct
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

//ToHex is a utility function that we will use to cast our nonce into a byte
func toHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
