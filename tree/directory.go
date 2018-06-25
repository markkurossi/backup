//
// directory.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"github.com/markkurossi/backup/storage"
)

type Directory struct {
	Type    Type
	Version Version
	Entries []DirectoryEntry
}

func (d *Directory) Serialize() ([]byte, error) {
	return Marshal(d)
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

func (d *Directory) Add(name string, mode uint32, modTime int64,
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
		Type:    TypeDirectory,
		Version: 1,
	}
}

type DirectoryEntry struct {
	Name    string
	Mode    uint32
	ModTime int64
	Entry   *storage.ID
}
