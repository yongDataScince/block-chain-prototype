package main

import (
	"fmt"

	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/chain"
)

func main() {
	c := chain.InitChain();
	c.AddBlock("First Block")
	c.AddBlock("Second Block")
	c.AddBlock("Thirth Block")
	c.AddBlock("Forth Block")
	c.AddBlock("Fifth Block")

	for i, block := range(c.Blocks) {
		fmt.Printf("Block: %x\n", i)
		fmt.Printf("Block data: %s\n", block.Data);
		fmt.Printf("Block Hash: %x\n", block.Hash);
		fmt.Println("Block Nonse: ", block.Nonse);

		fmt.Print("\n")
	}
}
