//
// element.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

type Element interface {
	Serialize() ([]byte, error)
}
