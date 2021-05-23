package chain

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
	b "github.com/kirillNovoseletskii/block-chain-prototype/pkg/block"
)

const (
	dbUrl = "./tmp/blocks"
)

type Chain struct {
	DataBase *badger.DB
	LastHash []byte
}

type ChainIterator struct {
	CurrentHash []byte
	DataBase 		*badger.DB
}

func (c *Chain) Iterator() *ChainIterator {
	return &ChainIterator{c.LastHash, c.DataBase}
}

func (iter *ChainIterator) Next() *b.Block {
	var encBlock []byte
	err := iter.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			log.Fatal(err)
		}

		err = item.Value(func(val []byte) error {
			encBlock = val
			return nil
		})

		return err
	})

	if err != nil {
		log.Fatal(err)
	}
	decBlock := b.Deserialize(encBlock);
	iter.CurrentHash = decBlock.PrevHash

	return decBlock
}

func (c *Chain) AddBlock(data string) {
	var lastHash []byte
	
	err := c.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"));
		if err != nil {
			log.Fatal(err)
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return err
		})

		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	newBlock := b.NewBlock(data, lastHash)

	err = c.DataBase.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Fatal(err)
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)
		c.LastHash = newBlock.Hash

		return err
	})

	if err != nil {
		log.Fatal(err)
	}
}

func genesis() *b.Block {
	return b.NewBlock("genesis", []byte{})
}

func InitChain() *Chain {
	var lastHash []byte;

	opts := badger.DefaultOptions(dbUrl)
	db, err := badger.Open(opts)

	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("lh")); if err == badger.ErrKeyNotFound {
			fmt.Println("Existing BlockChain not fount")
			gen := genesis()
			fmt.Println("Genesis Provided")
			err := txn.Set(gen.Hash, gen.Serialize());
			if err != nil {
				log.Fatal(err)
			}
			err = txn.Set([]byte("lh"), gen.Hash);
			if err != nil {
				log.Fatal(err)
			}
			lastHash = gen.Hash
			return nil
		} else {
			item, err := txn.Get([]byte("lh"));
			if err != nil {
				log.Fatal(err)
			}
			var value []byte
			err = item.Value(func(val []byte) error {
				value = val
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			lastHash = value
		}
		return err
	})
	return &Chain{
		DataBase: db,
		LastHash: lastHash,
	}
}
