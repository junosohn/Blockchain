package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	PrevHash   []byte
	Generation uint64
	Difficulty uint8
	Data       string
	Proof      uint64
	Hash       []byte
}

// Create new initial (generation 0) block.
func Initial(difficulty uint8) Block {
	b := new(Block)

	b.PrevHash = make([]byte, 32)
	b.Generation = 0
	b.Difficulty = difficulty
	b.Data = ""
	b.Proof = 0
	b.Hash = nil

	return *b
}

// Create new block to follow this block, with provided data.
func (prev_block Block) Next(data string) Block {
	b2 := new(Block)

	b2.PrevHash = prev_block.Hash
	b2.Generation = prev_block.Generation + 1
	b2.Difficulty = prev_block.Difficulty
	b2.Data = data
	b2.Proof = 0
	b2.Hash = nil

	return *b2
}

// Calculate the block's hash.
func (blk Block) CalcHash() []byte {
	prevhash := hex.EncodeToString(blk.PrevHash) + ":"
	generation := fmt.Sprintf("%v", blk.Generation) + ":"
	difficulty := fmt.Sprintf("%v", blk.Difficulty) + ":"
	data := blk.Data + ":"
	proof := fmt.Sprintf("%v", blk.Proof)

	// concatenate
	str := prevhash + generation + difficulty + data + proof

	//get SHA256 checksum
	hash := sha256.Sum256([]byte(str)[:])
	return hash[:]
}

// Is this block's hash valid?
func (blk Block) ValidHash() bool {
	nBytes := int(blk.Difficulty / 8)
	nBits := uint(blk.Difficulty % 8)
	ret := false
	t := blk.Hash


	if nBytes == 0 {
		// check only last byte
		checkLast := t[len(t)-nBytes-1]
		if checkLast % (1<<nBits) == 0 && checkLast != 0{
			// turn into binary
			bin := fmt.Sprintf("%b", checkLast)

			count := 0
			for i := 0; i < int(nBits); i++ {
				if bin[len(bin)-1] == '0' {
					count++
					bin = bin[0:len(bin)-1]
				}
			}

			// number of zero bits are EXACT
			if count == int(nBits) {
				ret = true
			}
		} else if checkLast % (1<<nBits) == 0 && checkLast == 0 {
			// 0 % (1<<nBits) will always equal 0, so return true
			ret = true
		}

	} else {
		// check all 'nBytes' bytes
		for i := 0; i < nBytes; i++ {

			checkLast := t[len(t)-i-1]

			if checkLast != 0 {
				return false
			}
		}

		// then check remaining nBits
		checkLast := t[len(t)-nBytes-1]

		if checkLast % (1<<nBits) == 0 && checkLast != 0{

			// turn into binary
			bin := fmt.Sprintf("%b", checkLast)

			count := 0
			for i := 0; i < int(nBits); i++ {
				if bin[len(bin)-1] == '0' {
					count++
					bin = bin[0:len(bin)-1]
				}
			}

			// number of zero bits are EXACT
			if count == int(nBits) {
				ret = true
			}
		} else if checkLast % (1<<nBits) == 0 && checkLast == 0 {
			// 0 % (1<<nBits) will always equal 0, so return true
			ret = true
		}

	}

	return ret
}

// Set the proof-of-work and calculate the block's "true" hash.
func (blk *Block) SetProof(proof uint64) {
	blk.Proof = proof
	blk.Hash = blk.CalcHash()
}
