package blockchain

import (
	"bytes"
)

type Blockchain struct {
	Chain []Block
}

func (chain *Blockchain) Add(blk Block) {
	// You can remove the panic() here if you wish.
	if !blk.ValidHash() {
		panic("adding block with invalid hash")
	}

	// append new block to the chain
	chain.Chain = append(chain.Chain, blk)
}

func (chain Blockchain) IsValid() bool {
	//	The initial block has previous hash all null bytes and is generation zero.
	for i := 0; i < len(chain.Chain[0].PrevHash); i++ {
		if chain.Chain[0].PrevHash[i] != 0 {
			return false
		}
		if chain.Chain[0].Generation != 0 {
			return false
		}
	}

	// Checking all other blocks in the chain
	for i := 0; i < (len(chain.Chain)-1); i++ {
		//	Each block has the same difficulty value.
		if chain.Chain[i].Difficulty != chain.Chain[i+1].Difficulty {
			return false
		}

		//	Each block has a generation value that is one more than the previous block.
		if chain.Chain[i].Generation != chain.Chain[i+1].Generation-1 {
			return false
		}

		//	Each block's previous hash matches the previous block's hash.
		if !bytes.Equal(chain.Chain[i+1].PrevHash, chain.Chain[i].Hash) {
			return false
		}

		//	Each block's hash value actually matches its contents.
		if !bytes.Equal(chain.Chain[i].CalcHash(), chain.Chain[i].Hash) {
			return false
		}

		//	Each block's hash value ends in difficulty null bits.
		if chain.Chain[i].ValidHash() == false {
			return false
		}

	}

	return true
}