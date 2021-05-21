package chain

import (
	b "github.com/kirillNovoseletskii/block-chain-prototype/pkg/block"
)

type Chain struct {
	Blocks []*b.Block
}

func (c *Chain) AddBlock(data string) {
	prevBlock := c.Blocks[len(c.Blocks)-1]
	newBlock := b.NewBlock(data, prevBlock.PrevHash)

	c.Blocks = append(c.Blocks, newBlock)
}

func genesis() *b.Block {
	return b.NewBlock("genesis", []byte{})
}

func InitChain() *Chain {
	return &Chain{
		Blocks: []*b.Block{genesis()},
	}
}
