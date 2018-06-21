//
// file.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io/ioutil"
	"os"

	"github.com/markkurossi/backup/tree"
)

type File struct {
	root string
	h    hash.Hash
}

func NewFile(root string) (Writer, error) {
	fmt.Printf("Initializing filesystem storage to '%s'\n", root)

	fileInfo, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("Filesystem storage root is not a directory")
	}

	return &File{
		root: root,
		h:    sha256.New(),
	}, nil
}

func (f *File) Write(data []byte) (*tree.ID, error) {
	f.h.Reset()
	f.h.Write(data)

	id := tree.NewID(f.h.Sum(nil))
	err := f.makeDirTree(id)
	if err != nil {
		return nil, err
	}
	path, err := f.makeFilename(id)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (f *File) makeFilename(id *tree.ID) (string, error) {
	if len(id.Data) < 2 {
		return "", fmt.Errorf("Invalid ID: %s", id)
	}
	return fmt.Sprintf("%s/%x/%x/%x",
		f.root, id.Data[:1], id.Data[1:2], id.Data[2:]), nil
}

func (f *File) makeDirTree(id *tree.ID) error {
	if len(id.Data) < 2 {
		return fmt.Errorf("Invalid ID: %s", id)
	}
	path := fmt.Sprintf("%s/%x/%x", f.root, id.Data[:1], id.Data[1:2])

	return os.MkdirAll(path, 0755)
}
