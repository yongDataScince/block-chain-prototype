package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	b "github.com/kirillNovoseletskii/block-chain-prototype/pkg/block"
	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/chain"
	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/handle"
)

// struct for cli comands
type CommandLine struct {
	blockChain *chain.Chain
}

// print blockchain usage
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tPrint Chain -> 'print'")
	fmt.Println("\tadd Block   -> 'add -block BLOCK_DATA' ")
	fmt.Println("\tclear db   -> 'clear' ")
	fmt.Println("\tExample: '$./main print'")
}

// get commands arguments
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// cli command for add block to chain
func (cli *CommandLine) addBlock(data string) {
	cli.blockChain.AddBlock(data)
	fmt.Println("Block added✅")
}

// cli command for print all blocks
func (cli *CommandLine) printChain() {
	fmt.Println()
	fmt.Println("Chain: ")
	iter := cli.blockChain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Block data: %s\n", block.Data)
		fmt.Printf("Block hash: %x\n", block.Hash)
		t, _ := strconv.ParseInt(fmt.Sprint(block.TimeStamp), 10, 64)
		fmt.Println("Block TimeStamp: ", time.Unix(t, 0))
		fmt.Println("Block Nonse: ", block.Nonse)
		pow := b.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// start cli
func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		handle.HandleError(err)
	case "print":
		err := printCmd.Parse(os.Args[2:])
		handle.HandleError(err)
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
