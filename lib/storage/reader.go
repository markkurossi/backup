//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package storage

// Reader reads data from the storage.
type Reader interface {
	// Read reads the data of the object id.
	Read(id ID) ([]byte, error)
}
