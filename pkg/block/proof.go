package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/handle"
)

const (
	Difficulty = 21
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (pow *ProofOfWork) InitData(nonse int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			toHex(pow.Block.TimeStamp),
			toHex(int64(nonse)),
			toHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// generate hash for block
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonse := 0

	for nonse < math.MaxInt64 {
		data := pow.InitData(nonse)
		hash = sha256.Sum256(data)

		fmt.Printf("\rhash selection: %x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonse++
		}
	}

	fmt.Println()
	fmt.Println()

	return nonse, hash[:]
}

// check block for PoW
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonse)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// function which create byte from number
func toHex(n int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, n)
	handle.HandleError(err)

	return buff.Bytes()
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	return &ProofOfWork{b, target}
}
