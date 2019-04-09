// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package cgtdb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"

	log "github.com/Sirupsen/logrus"
)

type LDBDatabase struct {
	fn string      // filename for reporting
	db *leveldb.DB // LevelDB instance
}

// NewLDBDatabase returns a LevelDB wrapped object.
func NewLDBDatabase(file string) (*LDBDatabase, error) {

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(file, nil)
	if err != nil {
		log.Warn("Open File Failed")
		return nil, err
	}
	return &LDBDatabase{
		fn:  file,
		db:  db,
	}, nil
}

// Path returns the path to the database directory.
func (db *LDBDatabase) Path() string {
	return db.fn
}

// Put puts the given key / value to the queue
func (db *LDBDatabase) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *LDBDatabase) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get returns the given key if it's present.
func (db *LDBDatabase) Get(key []byte) ([]byte, error) {
	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// Delete deletes the key from the queue and database
func (db *LDBDatabase) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *LDBDatabase) NewIterator() iterator.Iterator {
	return db.db.NewIterator(nil, nil)
}


func (db *LDBDatabase) Close() {
	err := db.db.Close()
	if err == nil {
		log.Info("Database closed")
	} else {
		log.Error("Failed to close database", "err", err)
	}
}

func (db *LDBDatabase) LDB() *leveldb.DB {
	return db.db
}
