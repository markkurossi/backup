//
// file.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"io"
)

type File interface {
	Size() int64
	Reader() io.Reader
}
