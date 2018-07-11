//
// directory.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"os"

	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/storage"
)

type Directory struct {
	ElementHeader
	Entries []DirectoryEntry
}

func (d *Directory) Serialize() ([]byte, error) {
	return encoding.Marshal(d)
}

func (d *Directory) IsDir() bool {
	return true
}

func (d *Directory) Directory() *Directory {
	return d
}

func (d *Directory) File() File {
	panic("Directory can't be converted to File")
}

func (d *Directory) Add(name string, mode os.FileMode, modTime int64,
	entry *storage.ID) {
	d.Entries = append(d.Entries, DirectoryEntry{
		Name:    name,
		Mode:    mode,
		ModTime: modTime,
		Entry:   entry,
	})
}

func NewDirectory() *Directory {
	return &Directory{
		ElementHeader: ElementHeader{
			Type:    TypeDirectory,
			Version: 1,
		},
	}
}

type DirectoryEntry struct {
	Name    string
	Mode    os.FileMode
	ModTime int64
	Entry   *storage.ID
}
