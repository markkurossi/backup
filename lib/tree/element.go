//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/storage"
)

// ElementHeader defines a common header for all element objects.
type ElementHeader struct {
	st      storage.Accessor `backup:"-"`
	Type    Type
	Version Version
}

// SetStorage sets the storage reference for the element.
func (hdr *ElementHeader) SetStorage(st storage.Accessor) {
	hdr.st = st
}

// Element defines an interface for storage elements.
type Element interface {
	SetStorage(st storage.Accessor)
	Serialize() ([]byte, error)
	IsDir() bool
	Directory() *Directory
	File() File
}

// DeserializeID deserializes an element id from the storage.
func DeserializeID(id storage.ID, st storage.Accessor) (Element, error) {
	if id.Undefined() {
		return nil, errors.New("Undefined ID")
	}
	data, err := st.Read(id)
	if err != nil {
		return nil, err
	}
	return Deserialize(data, st)
}

// Deserialize an element from the data.
func Deserialize(data []byte, st storage.Accessor) (Element, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("truncated element data")
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
		return nil, fmt.Errorf("unsupported tree element type %s", elementType)
	}

	err := encoding.Unmarshal(bytes.NewReader(data), element)
	if err != nil {
		return nil, err
	}

	element.SetStorage(st)

	return element, nil
}
