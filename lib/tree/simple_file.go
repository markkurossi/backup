//
// simple_file.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"io"

	"github.com/markkurossi/backup/lib/encoding"
)

type SimpleFile struct {
	ElementHeader
	Content []byte
}

func (f *SimpleFile) Serialize() ([]byte, error) {
	return encoding.Marshal(f)
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

func (f *SimpleFile) Reader() io.Reader {
	return &simpleReader{
		data: f.Content,
	}
}

type simpleReader struct {
	data []byte
}

func (r *simpleReader) Read(p []byte) (int, error) {
	read := copy(p, r.data)
	if read == 0 {
		return 0, io.EOF
	}
	r.data = r.data[read:]
	return read, nil
}

func NewSimpleFile(content []byte) *SimpleFile {
	return &SimpleFile{
		ElementHeader: ElementHeader{
			Type:    TypeSimpleFile,
			Version: 1,
		},
		Content: content,
	}
}
