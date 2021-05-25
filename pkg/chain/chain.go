package chain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	badger "githucom/dgraph-io/badger/v3"
	"githucom/kirillNovoseletskii/block-chain-prototype/pkg/handle"
)

const (
	dbUrl       = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from genesis"
)

type Chain struct {
	DataBase *badger.DB
	LastHash []byte // last hash of block
}

// type for printing chain
type ChainIterator struct {
	CurrentHash []byte
	DataBase    *badger.DB
}

func DBexist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (c *Chain) Iterator() *ChainIterator {
	return &ChainIterator{c.LastHash, c.DataBase}
}

// type for get next value of database
func (iter *ChainIterator) Next() *Block {
	var encBlock []byte
	err := iter.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		handle.HandleError(err)

		err = item.Value(func(val []byte) error {
			encBlock = val
			return nil
		})

		return err
	})

	handle.HandleError(err)

	decBlock := Deserialize(encBlock) // decode block from database
	iter.CurrentHash = decBlock.PrevHash

	return decBlock
}

// add blocks to chain and database
func (c *Chain) AddBlock(txs []*Transaction) {
	var lastHash []byte

	err := c.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		handle.HandleError(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return err
		})
		handle.HandleError(err)

		return err
	})
	handle.HandleError(err)

	newBlock := NewBlock(txs, lastHash)

	err = c.DataBase.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		handle.HandleError(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		handle.HandleError(err)
		c.LastHash = newBlock.Hash

		return err
	})

	handle.HandleError(err)
}

func (c *Chain) FindUnspentTx(addr string) []Transaction {
	var unspentTx []Transaction

	spentTxOs := make(map[string][]int)
	iter := c.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)

		Outputs:
			for outID, out := range tx.TxOutputs {
				if spentTxOs[txId] != nil {
					for _, spentOut := range spentTxOs[txId] {
						if spentOut == outID {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(addr) {
					unspentTx = append(unspentTx, *tx)
				}
				if tx.IsCoinbase() == false {
					for _, in := range tx.TxInputs {
						if in.CanUnlock(addr) {
							inTxId := hex.EncodeToString(in.ID)
							spentTxOs[inTxId] = append(spentTxOs[inTxId], in.Out)
						}
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTx
}

func (c *Chain) FindUTXO(addr string) []TxOutput {
	var UTXOs []TxOutput
	unsp := c.FindUnspentTx(addr)

	for _, tx := range unsp {
		for _, out := range tx.TxOutputs {
			if out.CanBeUnlocked(addr) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (c *Chain) FindSpentOuts(addr string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTx := c.FindUnspentTx(addr)
	accum := 0

	Work:
	for _, tx := range unspentTx {
		txId := hex.EncodeToString(tx.ID)
		for outId, out := range tx.TxOutputs {
			if out.CanBeUnlocked(addr) && accum < amount {
				accum += out.Value
				unspentOuts[txId] = append(unspentOuts[txId], outId)

				if accum >= amount {
					break Work
				}
			}
		}
	}

	return accum, unspentOuts
}

// create first(genesis) block of chain
func genesis(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// create chain
func InitChain(addr string) *Chain {
	var lastHash []byte

	if DBexist() {
		fmt.Println("Blockchain already exist")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbUrl)
	opts.Logger = nil
	db, err := badger.Open(opts)

	handle.HandleError(err)

	err = dUpdate(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(addr, genesisData)
		gen := genesis(cbtx)
		fmt.Println("Genesis Created")

		err := txn.Set(gen.Hash, gen.Serialize())
		handle.HandleError(err)

		err = txn.Set([]byte("lh"), gen.Hash)
		handle.HandleError(err)

		lastHash = gen.Hash

		return err
	})
	handle.HandleError(err)
	return &Chain{
		DataBase: db,
		LastHash: lastHash,
	}
}

func ContinueChain(addr string) *Chain {
	if !DBexist() {
		fmt.Println("Blockchain not exist")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbUrl)
	opts.Logger = nil

	db, err := badger.Open(opts)
	handle.HandleError(err)
	err = dUpdate(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		handle.HandleError(err)
		err = item.Value(func(val []byte) error { lastHash = val; return nil })
		handle.HandleError(err)
		return err
	})
	handle.HandleError(err)

	chain := Chain{db, lastHash}
	return &chain
}
