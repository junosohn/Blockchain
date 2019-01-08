package blockchain

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Testing CalcHash
func TestCalcHash(t *testing.T) {
	b0 := Initial(16)
	b0.Mine(1)
	b0CalcHash := hex.EncodeToString(b0.CalcHash())
	assert.Equal(t, b0CalcHash, "6c71ff02a08a22309b7dbbcee45d291d4ce955caa32031c50d941e3e9dbd0000")

	b1 := b0.Next("message")
	b1.Mine(1)
	b1CalcHash := hex.EncodeToString(b1.CalcHash())
	assert.Equal(t, b1CalcHash, "9b4417b36afa6d31c728eed7abc14dd84468fdb055d8f3cbe308b0179df40000")
}

// Testing ValidHash
func TestValidHash(t *testing.T) {
	b0 := Initial(19)
	b0.SetProof(87745)
	b1 := b0.Next("hash example 1234")

	b1.SetProof(1407891)
	assert.Equal(t, b0.ValidHash(), true)

	b1.SetProof(346082)
	assert.Equal(t, b1.ValidHash(), false)
}

// Testing Mining with difficulty=7
func TestMiningDifficulty7(t *testing.T) {
	b0 := Initial(7)
	b0.Mine(1)
	assert.Equal(t, b0.Proof, uint64(385))
	assert.Equal(t, hex.EncodeToString(b0.Hash), "379bf2fb1a558872f09442a45e300e72f00f03f2c6f4dd29971f67ea4f3d5300")

	b1 := b0.Next("this is an interesting message")
	b1.Mine(1)
	assert.Equal(t, b1.Proof, uint64(20))
	assert.Equal(t, hex.EncodeToString(b1.Hash), "4a1c722d8021346fa2f440d7f0bbaa585e632f68fd20fed812fc944613b92500")

	b2 := b1.Next("this is not interesting")
	b2.Mine(1)
	assert.Equal(t, b2.Proof, uint64(40))
	assert.Equal(t, hex.EncodeToString(b2.Hash), "ba2f9bf0f9ec629db726f1a5fe7312eb76270459e3f5bfdc4e213df9e47cd380")
}

// Testing Mining with difficulty=20
func TestMiningDifficulty20(t *testing.T) {
	b0 := Initial(20)
	b0.Mine(1)
	assert.Equal(t, b0.Proof, uint64(1209938))
	assert.Equal(t, hex.EncodeToString(b0.Hash), "19e2d3b3f0e2ebda3891979d76f957a5d51e1ba0b43f4296d8fb37c470600000")

	b1 := b0.Next("this is an interesting message")
	b1.Mine(1)
	assert.Equal(t, b1.Proof, uint64(989099))
	assert.Equal(t, hex.EncodeToString(b1.Hash), "a42b7e319ee2dee845f1eb842c31dac60a94c04432319638ec1b9f989d000000")

	b2 := b1.Next("this is not interesting")
	b2.Mine(1)
	assert.Equal(t, b2.Proof, uint64(1017262))
	assert.Equal(t, hex.EncodeToString(b2.Hash), "6c589f7a3d2df217fdb39cd969006bc8651a0a3251ffb50470cbc9a0e4d00000")
}

// Testing IsValid on Blockchain
func TestIsValid(t *testing.T) {
	b0 := Initial(10)
	b0.Mine(5)
	b1 := b0.Next("this is an interesting message")
	b1.Mine(5)
	b2 := b1.Next("this is not interesting")
	b2.Mine(5)

	blkChain := new(Blockchain)
	blkChain.Add(b0); blkChain.Add(b1); blkChain.Add(b2)

	if !blkChain.IsValid() {
		t.Error("Invalid Blockchain")
	}
}

// Additional testing on IsValid
func TestIsValidExtra(t *testing.T) {
	b0 := Initial(10)
	b0.Mine(5)
	b1 := b0.Next("this is an interesting message")
	b1.Mine(5)
	b2 := b1.Next("this is not interesting")
	b2.Mine(5)

	// Change difficulty value
	b0.Difficulty = 1
	blkChain2 := new(Blockchain)
	blkChain2.Add(b0); blkChain2.Add(b1); blkChain2.Add(b2)

	if !blkChain2.IsValid() {
		t.Error("Invalid Blockchain - blocks have different difficulty value")
	}

	// Change hash value
	b2.Hash[25] += 1
	blkChain3 := new(Blockchain)
	blkChain3.Add(b0); blkChain3.Add(b1); blkChain3.Add(b2)

	if !blkChain3.IsValid() {
		t.Error("Invalid Blockchain - block's hash value doesn't match its contents")
	}
}