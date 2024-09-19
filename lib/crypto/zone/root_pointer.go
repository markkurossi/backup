//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package zone

import (
	"github.com/markkurossi/backup/lib/storage"
)

// RootPointer implements zone root pointer.
type RootPointer struct {
	Version   byte
	Timestamp int64
	Pointer   storage.ID
	Digest    []byte
}
