//
// id.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type ID struct {
	Data []byte
}

func NewID(data []byte) ID {
	return ID{
		Data: data,
	}
}

func IDFromString(input string) (id ID, err error) {
	data, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	return NewID(data), nil
}

func (id ID) Undefined() bool {
	return len(id.Data) == 0
}

func (id ID) Equal(o ID) bool {
	return bytes.Equal(id.Data, o.Data)
}

func (id ID) String() string {
	if len(id.Data) > 16 {
		return fmt.Sprintf("%x...%x", id.Data[0:8], id.Data[len(id.Data)-8:])
	}
	return fmt.Sprintf("%x", id.Data)
}

func (id ID) ToFullString() string {
	return fmt.Sprintf("%x", id.Data)
}
