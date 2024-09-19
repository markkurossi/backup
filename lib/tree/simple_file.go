//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"io"

	"github.com/markkurossi/backup/lib/encoding"
)

// SimpleFile implements simple file objects.
type SimpleFile struct {
	ElementHeader
	Content []byte
}

// Serialize implements Element.Serialize.
func (f *SimpleFile) Serialize() ([]byte, error) {
	return encoding.Marshal(f)
}

// IsDir implements Element.IsDir.
func (f *SimpleFile) IsDir() bool {
	return false
}

// Directory implements Element.Directory.
func (f *SimpleFile) Directory() *Directory {
	panic("SimpleFile can't be converted to Directory")
}

// File implements Element.File.
func (f *SimpleFile) File() File {
	return f
}

// Size implements File.Size.
func (f *SimpleFile) Size() int64 {
	return int64(len(f.Content))
}

// Reader implements File.Reader.
func (f *SimpleFile) Reader() io.Reader {
	return &SimpleReader{
		data: f.Content,
	}
}

// SimpleReader implements io.Reader for simple file.
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

// Size returns the file size.
func (r *SimpleReader) Size() int64 {
	return int64(len(r.data))
}

// NewSimpleFile creates a new simple file object.
func NewSimpleFile(content []byte) *SimpleFile {
	return &SimpleFile{
		ElementHeader: ElementHeader{
			Type:    TypeSimpleFile,
			Version: 1,
		},
		Content: content,
	}
}
