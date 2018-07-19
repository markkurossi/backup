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

var typeNames = map[Type]string{
	TypeSimpleFile:  "simple-file",
	TypeChunkedFile: "chunked-file",
	TypeDirectory:   "directory",
	TypeSnapshot:    "snapshot",
}

func (t Type) String() string {
	name, ok := typeNames[t]
	if ok {
		return name
	}
	return fmt.Sprintf("{Type %d}", t)
}

const (
	TypeSimpleFile Type = iota + 1
	TypeChunkedFile
	TypeDirectory
	TypeSnapshot
)

type Version uint8
