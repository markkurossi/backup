//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package persistence

// Flags define flags for persistence operations.
type Flags uint

const (
	// NoCache forces reading from the storage bypassing any cached
	// content.
	NoCache Flags = 1 << iota
)

// Reader defines persistence reader interface.
type Reader interface {
	// Exists tests if the specified key exists in the namespace.
	Exists(namespace, key string) (bool, error)

	// Get gets the data of the specified key in the namespace.
	Get(namespace, key string, flags Flags) ([]byte, error)

	// GetAll returns all keys and their values from the namespae.
	GetAll(namespace string) (map[string][]byte, error)
}
