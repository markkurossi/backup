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
	return fmt.Sprintf("%x", id.Data)
}
