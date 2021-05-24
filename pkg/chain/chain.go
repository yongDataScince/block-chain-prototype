package chain

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
	b "github.com/kirillNovoseletskii/block-chain-prototype/pkg/block"
	"github.com/kirillNovoseletskii/block-chain-prototype/pkg/handle"
)

const (
	dbUrl = "./tmp/blocks"
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

func (c *Chain) Iterator() *ChainIterator {
	return &ChainIterator{c.LastHash, c.DataBase}
}

// type for get next value of database
func (iter *ChainIterator) Next() *b.Block {
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

	decBlock := b.Deserialize(encBlock) // decode block from database
	iter.CurrentHash = decBlock.PrevHash

	return decBlock
}

// add blocks to chain and database
func (c *Chain) AddBlock(data string) {
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

	newBlock := b.NewBlock(data, lastHash)

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

// create first(genesis) block of chain
func genesis() *b.Block {
	return b.NewBlock("genesis", []byte{})
}

// create chain
func InitChain() *Chain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbUrl)
	opts.Logger = nil
	db, err := badger.Open(opts)

	handle.HandleError(err)

	db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("lh"))
		if err == badger.ErrKeyNotFound {
			// check if block chain not found in database
			fmt.Println("Existing BlockChain not fount")
			gen := genesis()
			fmt.Println("Genesis Provided")
			err := txn.Set(gen.Hash, gen.Serialize()) // write genesis block in database
			handle.HandleError(err)

			err = txn.Set([]byte("lh"), gen.Hash)
			handle.HandleError(err)

			lastHash = gen.Hash
			return nil
		} else {
			item, err := txn.Get([]byte("lh"))
			handle.HandleError(err)

			var value []byte
			err = item.Value(func(val []byte) error {
				value = val
				return nil
			})
			handle.HandleError(err)
			lastHash = value
		}
		return err
	})
	return &Chain{
		DataBase: db,
		LastHash: lastHash,
	}
}
