//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"io"
)

// File implements interface for files.
type File interface {
	Size() int64
	Reader() io.Reader
}
