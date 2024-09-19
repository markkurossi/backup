//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"io"

	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/storage"
)

// ChunkedFile implements chunked file objects.
type ChunkedFile struct {
	ElementHeader
	ContentSize int64
	Chunks      []Chunk
}

// Serialize implements Element.Serialize.
func (c *ChunkedFile) Serialize() ([]byte, error) {
	return encoding.Marshal(c)
}

// IsDir implements Element.IsDir.
func (c *ChunkedFile) IsDir() bool {
	return false
}

// Directory implements Element.Directory.
func (c *ChunkedFile) Directory() *Directory {
	panic("ChunkedFile can't be converted to Directory")
}

// File implements Element.File.
func (c *ChunkedFile) File() File {
	return c
}

// Size implements File.Size.
func (c *ChunkedFile) Size() int64 {
	return c.ContentSize
}

// Reader implements File.Reader.
func (c *ChunkedFile) Reader() io.Reader {
	return &chunkReader{
		st:     c.st,
		chunks: c.Chunks,
	}
}

type chunkReader struct {
	st     storage.Accessor
	chunks []Chunk
	data   []byte
}

func (r *chunkReader) Read(p []byte) (n int, err error) {
	if len(r.data) == 0 {
		if len(r.chunks) == 0 {
			return 0, io.EOF
		}
		data, err := r.st.Read(r.chunks[0].Content)
		if err != nil {
			return 0, err
		}
		r.data = data
		r.chunks = r.chunks[1:]
	}

	read := copy(p, r.data)
	if read == 0 {
		return 0, io.EOF
	}
	r.data = r.data[read:]
	return read, nil
}

// Add adds a file chunk.
func (c *ChunkedFile) Add(size int64, chunk storage.ID) {
	c.Chunks = append(c.Chunks, Chunk{
		Size:    size,
		Content: chunk,
	})
}

// NewChunkedFile creates a new chunked file object.
func NewChunkedFile(size int64) *ChunkedFile {
	return &ChunkedFile{
		ElementHeader: ElementHeader{
			Type:    TypeChunkedFile,
			Version: 1,
		},
		ContentSize: size,
	}
}

// Chunk implements one file content chunk.
type Chunk struct {
	Size    int64
	Content storage.ID
}
