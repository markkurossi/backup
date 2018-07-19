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

	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/storage"
)

type ElementHeader struct {
	st      storage.Accessor `backup:"-"`
	Type    Type
	Version Version
}

func (hdr *ElementHeader) SetStorage(st storage.Accessor) {
	hdr.st = st
}

type Element interface {
	SetStorage(st storage.Accessor)
	Serialize() ([]byte, error)
	IsDir() bool
	Directory() *Directory
	File() File
}

func Deserialize(id storage.ID, st storage.Accessor) (Element, error) {
	data, err := st.Read(id)
	if err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, fmt.Errorf("Truncated element data")
	}
	var element Element

	elementType := Type(data[0])
	switch elementType {
	case TypeSimpleFile:
		element = new(SimpleFile)

	case TypeChunkedFile:
		element = new(ChunkedFile)

	case TypeDirectory:
		element = new(Directory)

	case TypeSnapshot:
		element = new(Snapshot)

	default:
		return nil, fmt.Errorf("Unsupported tree element type %s", elementType)
	}

	err = encoding.Unmarshal(bytes.NewReader(data), element)
	if err != nil {
		return nil, err
	}

	element.SetStorage(st)

	return element, nil
}
