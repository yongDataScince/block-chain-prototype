package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/handle"
)

// block struct
type Block struct {
	TimeStamp int64
	Hash      []byte
	PrevHash  []byte
	Transactions []*Transaction
	Nonse     int
}

// function for  encode block to `bytes` for add to database
func (b *Block) Serialize() []byte {
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	handle.HandleError(err)

	return res.Bytes()
}

// function for decode block from database
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	handle.HandleError(err)

	return &block
}

func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// create new block
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := &Block{
		TimeStamp: time.Now().Unix(),
		PrevHash:  prevHash,
		Transactions: txs,
		Nonse:     0,
	}

	pow := NewProof(b)
	nonse, hash := pow.Run() // Proof Of Work validate for block

	b.Hash = hash[:]
	b.Nonse = nonse

	return b
}
