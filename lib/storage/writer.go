//
// writer.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

type Writer interface {
	// Write writes the data to the storage. The function
	// returns the object ID.
	Write(data []byte) (ID, error)
}
