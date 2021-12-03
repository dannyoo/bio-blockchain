package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)
// returns public key from file location
func loadPublicKey(filename string) rsa.PublicKey {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode([]byte(content))
	key, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return *key
}

// load private key from file if it doesn't exist create it.
func loadPrivateKey() rsa.PrivateKey {
	if _, err := os.Stat("private.pem"); err == nil {
		content, err := os.ReadFile("private.pem")
		if err != nil {
			log.Fatal(err)
		}
		block, _ := pem.Decode([]byte(content))
		key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
		return *key
	} else {
		saveToFile()
		return loadPrivateKey()
	}
}

// saves private and public key to file
func saveToFile() {

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		ErrorHandle(err)
		os.Exit(1)
	}

	publickey := &privatekey.PublicKey

	pemfile, err := os.Create("private.pem")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var pemkey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}

	err = pem.Encode(pemfile, pemkey)

	if err != nil {
		ErrorHandle(err)
		os.Exit(1)
	}

	pemfile.Close()

	// save public key

	publicKeyBytes := x509.MarshalPKCS1PublicKey(publickey)

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create("public.pem")
	if err != nil {
		fmt.Printf("error when create public.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		fmt.Printf("error when encode public pem: %s \n", err)
		os.Exit(1)
	}

}

// encrypts string
func rsaEncrypt(secretMessage string, key rsa.PublicKey) string {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	ErrorHandle(err)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// decrypts string
func rsaDecrypt(cipherText string, privKey rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	return string(plaintext), err
}

