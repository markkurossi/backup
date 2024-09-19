//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package storage

// Accessor defines storage accessor interface.
type Accessor interface {
	Reader
	Writer
}
