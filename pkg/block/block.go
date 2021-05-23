package block

import (
	"time"
)

type Block struct {
	TimeStamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
	Nonse		  int
}

func NewBlock(data string, prevHash []byte) *Block {
	b := &Block{
		TimeStamp: time.Now().Unix(),
		PrevHash:  prevHash,
		Data:      []byte(data),
		Nonse: 0,
	}

	pow := NewProof(b)
	nonse, hash := pow.Run()

	b.Hash = hash[:]
	b.Nonse = nonse

	return b
}
