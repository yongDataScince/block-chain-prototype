package block

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	TimeStamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
	Nonse		  int
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res);
	err := encoder.Encode(b);
	if err != nil {
		log.Fatal(err)
	}

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block);
	if err != nil {
		log.Fatal(err)
	}

	return &block
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
