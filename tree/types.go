//
// types.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"fmt"
)

type Type uint8

func (t Type) String() string {
	switch t {
	case TypeFile:
		return "file"

	case TypeChunkedFile:
		return "chunked-file"

	case TypeDirectory:
		return "directory"

	default:
		return fmt.Sprintf("{Type %d}", t)
	}
}

const (
	TypeFile Type = iota
	TypeChunkedFile
	TypeDirectory
)
