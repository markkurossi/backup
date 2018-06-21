//
// types.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

type Type uint8

const (
	TypeFile Type = iota
	TypeChunkedFile
	TypeDirectory
)
