//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package persistence

// Writer defines persistence writer interface.
type Writer interface {
	// Set sets the data to the specified key in the namespace.
	Set(namespace, key string, data []byte) error
}
