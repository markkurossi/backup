//
// accessor.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package persistence

type Accessor interface {
	Reader
	Writer
}
