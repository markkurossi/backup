//
// element.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"bytes"
	"fmt"

	"github.com/markkurossi/backup/storage"
)

type Element interface {
	Serialize() ([]byte, error)
	IsDir() bool
	Directory() *Directory
	File() File
}

func Deserialize(id *storage.ID, reader storage.Reader) (Element, error) {
	data, err := reader.Read(id)
	if err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, fmt.Errorf("Truncated element data")
	}
	var element Element

	switch Type(data[0]) {
	case TypeSimpleFile:
		element = new(SimpleFile)

	case TypeChunkedFile:
		element = new(ChunkedFile)

	case TypeDirectory:
		element = new(Directory)

	default:
		return nil, fmt.Errorf("Unsupported tree element type %s")
	}

	err = Unmarshal(bytes.NewReader(data), element)
	if err != nil {
		return nil, err
	}

	return element, nil
}
