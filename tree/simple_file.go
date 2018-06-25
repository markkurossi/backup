//
// simple_file.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

type SimpleFile struct {
	Type    Type
	Version Version
	Content []byte
}

func (f *SimpleFile) Serialize() ([]byte, error) {
	return Marshal(f)
}

func (f *SimpleFile) IsDir() bool {
	return false
}

func (f *SimpleFile) Directory() *Directory {
	panic("SimpleFile can't be converted to Directory")
}

func (f *SimpleFile) File() File {
	return f
}

func (f *SimpleFile) Size() int64 {
	return int64(len(f.Content))
}

func NewSimpleFile(content []byte) *SimpleFile {
	return &SimpleFile{
		Type:    TypeSimpleFile,
		Version: 1,
		Content: content,
	}
}
