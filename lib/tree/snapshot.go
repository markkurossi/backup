//
// snapshot.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/storage"
)

type Snapshot struct {
	ElementHeader
	Timestamp int64
	Root      storage.ID
	Parent    storage.ID
}

func (s *Snapshot) Serialize() ([]byte, error) {
	return encoding.Marshal(s)
}

func (s *Snapshot) IsDir() bool {
	return false
}

func (s *Snapshot) Directory() *Directory {
	panic("Snapshot can't be converted to Directory")
}

func (s *Snapshot) File() File {
	panic("Snapshot can't be converted to File")
}

func NewSnapshot() *Snapshot {
	return &Snapshot{
		ElementHeader: ElementHeader{
			Type:    TypeSnapshot,
			Version: 1,
		},
	}
}
