//
// reader.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

import (
	"github.com/markkurossi/backup/tree"
)

type Reader interface {
	Read(id *tree.ID) ([]byte, error)
}
