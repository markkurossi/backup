//
// reader.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

type Reader interface {
	Read(id ID) ([]byte, error)
}
