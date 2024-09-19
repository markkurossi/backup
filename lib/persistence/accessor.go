//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package persistence

// Accessor defines persistence storage accessor.
type Accessor interface {
	Reader
	Writer
}

var (
	_ Accessor = &Filesystem{}
	_ Accessor = &HTTP{}
)
