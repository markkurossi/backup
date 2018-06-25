//
// accessor.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

type Accessor interface {
	Reader
	Writer
}
