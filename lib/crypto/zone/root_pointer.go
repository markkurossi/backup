//
// root_pointer.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package zone

import (
	"github.com/markkurossi/backup/lib/storage"
)

type RootPointer struct {
	Timestamp int64
	Pointer   storage.ID
	Digest    []byte
}
