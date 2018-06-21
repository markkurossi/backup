//
// directory.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

func NewDirectory() *Directory {
	return &Directory{
		Type: TypeDirectory,
	}
}

type Directory struct {
	Type    Type
	Entries []DirectoryEntry
}

func (d *Directory) Serialize() ([]byte, error) {
	return Marshal(d)
}

func (d *Directory) Add(name string, mode uint32, modTime int64, entry *ID) {
	d.Entries = append(d.Entries, DirectoryEntry{
		Name:    name,
		Mode:    mode,
		ModTime: modTime,
		Entry:   entry,
	})
}

type DirectoryEntry struct {
	Name    string
	Mode    uint32
	ModTime int64
	Entry   *ID
}
