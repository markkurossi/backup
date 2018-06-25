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

	"github.com/markkurossi/backup/tree"
)

type Null struct {
	h hash.Hash
}

func NewNull() *Null {
	return &Null{
		h: sha256.New(),
	}
}

func (n *Null) Write(data []byte) (*tree.ID, error) {
	n.h.Reset()
	n.h.Write(data)

	return tree.NewID(n.h.Sum(nil)), nil
}

func (n *Null) Read(id *tree.ID) ([]byte, error) {
	return nil, fmt.Errorf("Data not found")
}
