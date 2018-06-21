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
	Size    int64
	Content *ID
}

func (f *File) Serialize() ([]byte, error) {
	return Marshal(f)
}

func NewFile(size int64, content *ID) *File {
	return &File{
		Type:    TypeFile,
		Size:    size,
		Content: content,
	}
}
