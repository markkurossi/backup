//
// writer.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

import (
	"github.com/markkurossi/backup/tree"
)

type Writer interface {
	Write(data []byte) (*tree.ID, error)
}
