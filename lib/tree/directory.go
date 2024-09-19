//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"os"

	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/storage"
)

// Directory implements directory objects.
type Directory struct {
	ElementHeader
	Entries []DirectoryEntry
}

// Serialize implements Element.Serialize.
func (d *Directory) Serialize() ([]byte, error) {
	return encoding.Marshal(d)
}

// IsDir implements Element.IsDir.
func (d *Directory) IsDir() bool {
	return true
}

// Directory implements Element.Directory.
func (d *Directory) Directory() *Directory {
	return d
}

// File implements Element.File.
func (d *Directory) File() File {
	panic("Directory can't be converted to File")
}

// Add adds an entry to the directory.
func (d *Directory) Add(name string, mode os.FileMode, modTime int64,
	entry storage.ID) {
	d.Entries = append(d.Entries, DirectoryEntry{
		Name:    name,
		Mode:    mode,
		ModTime: modTime,
		Entry:   entry,
	})
}

// NewDirectory creates a new directory object.
func NewDirectory() *Directory {
	return &Directory{
		ElementHeader: ElementHeader{
			Type:    TypeDirectory,
			Version: 1,
		},
	}
}

// DirectoryEntry implements a directory entry.
type DirectoryEntry struct {
	Name    string
	Mode    os.FileMode
	ModTime int64
	Entry   storage.ID
}
