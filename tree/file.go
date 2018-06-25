//
// file.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

type File struct {
	Type    Type
	Version Version
	Content []byte
}

func (f *File) Serialize() ([]byte, error) {
	return Marshal(f)
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Directory() *Directory {
	panic("File can't be converted to Directory")
}

func (f *File) File() *File {
	return f
}

func NewFile(content []byte) *File {
	return &File{
		Type:    TypeFile,
		Version: 1,
		Content: content,
	}
}
