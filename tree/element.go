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
)

type Element interface {
	Serialize() ([]byte, error)
	IsDir() bool
	Directory() *Directory
	File() *File
}

func Deserialize(data []byte) (Element, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("Truncated element data")
	}
	var element Element

	switch Type(data[0]) {
	case TypeFile:
		element = new(File)

	case TypeDirectory:
		element = new(Directory)

	default:
		return nil, fmt.Errorf("Unsupported tree element type %s")
	}

	err := Unmarshal(bytes.NewReader(data), element)
	if err != nil {
		return nil, err
	}

	return element, nil
}
