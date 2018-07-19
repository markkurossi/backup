//
// null.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

import (
	"crypto/sha256"
	"fmt"
	"hash"
)

type Null struct {
	h hash.Hash
}

func NewNull() *Null {
	return &Null{
		h: sha256.New(),
	}
}

func (n *Null) Write(data []byte) (ID, error) {
	n.h.Reset()
	n.h.Write(data)

	return NewID(n.h.Sum(nil)), nil
}

func (n *Null) Read(id ID) ([]byte, error) {
	return nil, fmt.Errorf("Data not found")
}
