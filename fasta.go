package main

import (
	"bufio"
	"log"
	"os"
)

// Read fasta files
func readFasta(file string) (label, seq string) {
	// open file
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		txt := scanner.Text()
		if txt[0] == '>' {
			label = txt[1:]
		} else {
			seq = txt
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}
