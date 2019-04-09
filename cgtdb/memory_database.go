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
	"errors"
)

/*
 * This is a test memory database. Do not use for any production it does not get persisted
 */
type MemDatabase struct {
	db   map[string][]byte
}

func NewMemDatabase() *MemDatabase {
	return &MemDatabase{
		db: make(map[string][]byte),
	}
}

func (db *MemDatabase) Put(key []byte, value []byte) error {
	db.db[string(key)] = CopyBytes(value)
	return nil
}

func (db *MemDatabase) Has(key []byte) (bool, error) {
	_, ok := db.db[string(key)]
	return ok, nil
}

func (db *MemDatabase) Get(key []byte) ([]byte, error) {
	if entry, ok := db.db[string(key)]; ok {
		return CopyBytes(entry), nil
	}
	return nil, errors.New("not found")
}

func (db *MemDatabase) Keys() [][]byte {
	keys := [][]byte{}
	for key := range db.db {
		keys = append(keys, []byte(key))
	}
	return keys
}

func (db *MemDatabase) Delete(key []byte) error {
	delete(db.db, string(key))
	return nil
}

func (db *MemDatabase) Close() {}

func (db *MemDatabase) Len() int { return len(db.db) }

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes,b)

	return
}
