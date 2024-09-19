//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"fmt"
)

// Type defines the tree object types.
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

// Object tyhpes
const (
	TypeSimpleFile Type = iota + 1
	TypeChunkedFile
	TypeDirectory
	TypeSnapshot
)

// Version defines object version.
type Version uint8

// FileSize defines file size in bytes.
type FileSize int64

func (size FileSize) String() string {
	if size > 1024*1024 {
		return fmt.Sprintf("%d MB", size/(1024*1024))
	} else if size > 1024 {
		return fmt.Sprintf("%d kB", size/1024)
	} else {
		return fmt.Sprintf("%d B", size)
	}
}
