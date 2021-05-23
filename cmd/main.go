package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	b "github.com/kirillNovoseletskii/block-chain-prototype/pkg/block"
	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/chain"
)

type CommandLine struct {
	blockChain *chain.Chain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tPrint Chain -> 'print'")
	fmt.Println("\tadd Block -> 'add -block BLOCK_DATA' ")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockChain.AddBlock(data)
	fmt.Println("Block added✅")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockChain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Block data: %s\n", block.Data)
		fmt.Printf("Block hash: %x\n", block.Hash)
		pow := b.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "print":
		err := printCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)
	chain := chain.InitChain()
	defer chain.DataBase.Close()

	cli := CommandLine{chain}
	cli.run()
}
