//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"fmt"

	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/storage"
)

// Snapshot implements snapshot objects.
type Snapshot struct {
	ElementHeader
	Timestamp int64
	Size      FileSize
	Root      storage.ID
	Parent    storage.ID
}

func (s *Snapshot) String() string {
	return fmt.Sprintf("Snapshot %s (%s)", s.Root, s.Size)
}

// Serialize implements Element.Serialize.
func (s *Snapshot) Serialize() ([]byte, error) {
	return encoding.Marshal(s)
}

// IsDir implements Element.IsDir.
func (s *Snapshot) IsDir() bool {
	return false
}

// Directory implements Element.Directory.
func (s *Snapshot) Directory() *Directory {
	panic("Snapshot can't be converted to Directory")
}

// File implements Element.File.
func (s *Snapshot) File() File {
	panic("Snapshot can't be converted to File")
}

// NewSnapshot creates a new snapshot object.
func NewSnapshot() *Snapshot {
	return &Snapshot{
		ElementHeader: ElementHeader{
			Type:    TypeSnapshot,
			Version: 1,
		},
	}
}
