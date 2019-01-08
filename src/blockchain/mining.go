package blockchain

import (
	"work_queue"
)

type miningWorker struct {
	blk Block
	start uint64
	end uint64
}

// implement work_queue.Worker
func (mw miningWorker) Run() interface{} {
	mr := new(MiningResult)
	mr.Found = false

	// from start to end (inclusive)
	for i := mw.start; i <= mw.end; i++ {
		mw.blk.Proof = i
		mw.blk.Hash = mw.blk.CalcHash()

		if mw.blk.ValidHash() == true {
			mr.Proof = i
			mr.Found = true
			return *mr
		}
	}

	return *mr
}


type MiningResult struct {
	Proof uint64 // proof-of-work value, if found.
	Found bool   // true if valid proof-of-work was found.
}

// Mine the range of proof values, by breaking up into chunks and checking
// "workers" chunks concurrently in a work queue. Should return shortly after a result
// is found.
func (blk Block) MineRange(start uint64, end uint64, workers uint64, chunks uint64) MiningResult {
	wq := work_queue.Create(uint(workers), uint(chunks))
	//fmt.Printf("start=%v, end=%v, CHUNKS=%v\n", start, end, chunks)
	mr := new(MiningResult)
	mr.Found = false

	// check if chunk size is smaller than given number of chunks
	c := (end - start) / chunks
	if c == 0 {
		c = end
	}

	for i := start; i < end; i+=c {
		mw := new(miningWorker)
		mw.blk = blk
		mw.start = i

		// prevent overflow
		if i + c - 1 <= end {
			mw.end = i + c - 1
		} else {
			mw.end = end
		}

		wq.Enqueue(mw)
	}

	for {
		r := <- wq.Results
		if r.(MiningResult).Found == true {
			mr.Found = true
			mr.Proof = r.(MiningResult).Proof
			wq.Shutdown()
			break
		}
	}

	return *mr

}

// Call .MineRange with some reasonable values that will probably find a result.
// Good enough for testing at least. Updates the block's .Proof and .Hash if successful.
func (blk *Block) Mine(workers uint64) bool {
	reasonableRangeEnd := uint64(4 * 1 << blk.Difficulty) // 4 * 2^(bits that must be zero)
	mr := blk.MineRange(0, reasonableRangeEnd, workers, 128)
	if mr.Found {
		blk.SetProof(mr.Proof)
	}
	return mr.Found
}