package main

import (
	"fmt"

	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/chain"
)

// type CommandLine struct {
// 	blockChain *chain.Chain
// }

func main() {
	block_chain := chain.InitChain()

	iter := block_chain.Iterator();
	i := 0
	for i < 3 {
		block := iter.Next()
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
		i ++
	}
	
}
