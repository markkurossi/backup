//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package storage

// Writer writes data to the storage.
type Writer interface {
	// Write writes the data to the storage. The function
	// returns the object ID.
	Write(data []byte) (ID, error)
}
