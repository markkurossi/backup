//
// reader.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package persistence

type Reader interface {
	// Exists tests if the specified key exists in the namespace.
	Exists(namespace, key string) (bool, error)

	// Get gets the data of the specified key in the namespace.
	Get(namespace, key string) ([]byte, error)

	// GetAll returns all keys and their values from the namespae.
	GetAll(namespace string) (map[string][]byte, error)
}
