//
// id.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"encoding/hex"
	"fmt"
)

type ID struct {
	Data []byte
}

func NewID(data []byte) *ID {
	return &ID{
		Data: data,
	}
}

func IDFromString(id string) (*ID, error) {
	data, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}
	return NewID(data), nil
}

func (id *ID) String() string {
	if len(id.Data) > 16 {
		return fmt.Sprintf("%x...%x", id.Data[0:8], id.Data[len(id.Data)-8:])
	}
	return fmt.Sprintf("%x", id.Data)
}

func (id *ID) ToFullString() string {
	return fmt.Sprintf("%x", id.Data)
}
