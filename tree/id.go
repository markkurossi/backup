//
// id.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
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

func (id *ID) String() string {
	if len(id.Data) > 16 {
		return fmt.Sprintf("%x...%x", id.Data[0:8], id.Data[len(id.Data)-8:])
	}
	return fmt.Sprintf("%x", id.Data)
}
