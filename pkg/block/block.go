package block

import (
	"bytes"
	"strconv"
	"time"
)

type Block struct {
	TimeStamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
}

func (b *Block) driveHash() {
	timestamp := []byte(strconv.FormatInt(b.TimeStamp, 10))
	info := bytes.Join([][]byte{b.Data, b.Hash, timestamp}, []byte{})

	b.Hash = info[:]
}

func NewBlock(data string, prevHash []byte) *Block{
	b := &Block{
		TimeStamp: time.Now().Unix(),
		PrevHash: prevHash,
		Data: []byte(data),
	}

	b.driveHash()

	return b
}
