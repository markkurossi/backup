//
// writer.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package persistence

type Writer interface {
	// Set sets the data to the specified key in the namespace.
	Set(namespace, key string, data []byte) error
}
