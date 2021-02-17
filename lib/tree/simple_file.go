//
// simple_file.go
//
// Copyright (c) 2018-2021 Markku Rossi
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
	return &SimpleReader{
		data: f.Content,
	}
}

type SimpleReader struct {
	data []byte
	ofs  int
}

func (r *SimpleReader) Read(p []byte) (int, error) {
	read := copy(p, r.data[r.ofs:])
	if read == 0 {
		return 0, io.EOF
	}
	r.ofs += read
	return read, nil
}

func (r *SimpleReader) Size() int64 {
	return int64(len(r.data))
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
