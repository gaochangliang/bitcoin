package main

import (
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	handleError(err)
	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		handleError(err)
		err = b.Put([]byte("l"), newBlock.Hash)
		handleError(err)

		bc.tip = newBlock.Hash

		return nil
	})
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	handleError(err)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			handleError(err)
			err = b.Put(genesis.Hash, genesis.Serialize())
			handleError(err)
			err = b.Put([]byte("l"), genesis.Hash)
			handleError(err)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	handleError(err)
	bc := Blockchain{tip, db}

	return &bc
}