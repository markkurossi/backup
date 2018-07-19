//
// writer.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

type Writer interface {
	Write(data []byte) (ID, error)
}
