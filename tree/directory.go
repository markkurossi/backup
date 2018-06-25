//
// directory.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

type Directory struct {
	Type    Type
	Version Version
	Entries []DirectoryEntry
}

func (d *Directory) Serialize() ([]byte, error) {
	data, err := Marshal(d)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Directory) IsDir() bool {
	return true
}

func (d *Directory) Directory() *Directory {
	return d
}

func (d *Directory) File() *File {
	panic("Directory can't be converted to File")
}

func (d *Directory) Add(name string, mode uint32, modTime int64, entry *ID) {
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
	Entry   *ID
}
