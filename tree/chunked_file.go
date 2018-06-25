//
// chunked_file.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

type ChunkedFile struct {
	Type        Type
	Version     Version
	ContentSize int64
	Chunks      []Chunk
}

func (c *ChunkedFile) Serialize() ([]byte, error) {
	return Marshal(c)
}

func (c *ChunkedFile) IsDir() bool {
	return false
}

func (c *ChunkedFile) Directory() *Directory {
	panic("ChunkedFile can't be converted to Directory")
}

func (c *ChunkedFile) File() File {
	return c
}

func (c *ChunkedFile) Size() int64 {
	return c.ContentSize
}

func (c *ChunkedFile) Add(size int64, chunk *ID) {
	c.Chunks = append(c.Chunks, Chunk{
		Size:    size,
		Content: chunk,
	})
}

func NewChunkedFile(size int64) *ChunkedFile {
	return &ChunkedFile{
		Type:        TypeChunkedFile,
		Version:     1,
		ContentSize: size,
	}
}

type Chunk struct {
	Size    int64
	Content *ID
}
