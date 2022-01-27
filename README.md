# Bio Blockchain

## Getting Started: 

1. You need to run `go get github.com/dgraph-io/badger` before building. 
2. go build
3. `./blockchain` to see usage

### Quick Note
- Special Notes: Private and Public Key are generated saved as .pem files if not found in directory.

## Changes to "features and requirements":
- Personal id is replaced with public key.
- Did not complete decentralized networking interface for nodes on the network to communicate. It could be for future work if the project gains traction.
- Added encryption of data and public/private key generation which was not mentioned in the feature requirments.
- Added fasta file reading for the input of biological data when adding a block or transfering.
