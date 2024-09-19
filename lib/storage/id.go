//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package storage

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

var (
	// EmptyID is an undefined empty ID.
	EmptyID = ID{}
)

// ID defines the storage ID.
type ID struct {
	Data []byte
}

// NewID creates an ID from the data.
func NewID(data []byte) ID {
	return ID{
		Data: data,
	}
}

// IDFromString creates an ID from the input string.
func IDFromString(input string) (id ID, err error) {
	data, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	return NewID(data), nil
}

// Undefined tests if the ID is undefined.
func (id ID) Undefined() bool {
	return len(id.Data) == 0
}

// Equal tests if the argument ID is equal to this one.
func (id ID) Equal(o ID) bool {
	return bytes.Equal(id.Data, o.Data)
}

func (id ID) String() string {
	if len(id.Data) > 16 {
		return fmt.Sprintf("%x...%x", id.Data[0:8], id.Data[len(id.Data)-8:])
	}
	return fmt.Sprintf("%x", id.Data)
}

// ToFullString returns the full ID string.
func (id ID) ToFullString() string {
	return fmt.Sprintf("%x", id.Data)
}
