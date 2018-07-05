//
// key.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package identity

import (
	"bytes"
	"fmt"

	"github.com/markkurossi/backup/lib/encoding"
)

type Key interface {
	Marshal() ([]byte, error)
}

type KeyType int

func (t KeyType) String() string {
	switch t {
	case KeyRSA:
		return "rsa"

	default:
		return fmt.Sprintf("{KeyType %d}", t)
	}
}

const (
	KeyRSA KeyType = 0
)

type KeyData struct {
	Name string
	Type KeyType
	Data []byte
}

func Unmarshal(data []byte) (Key, error) {
	keyData := new(KeyData)
	if err := encoding.Unmarshal(bytes.NewReader(data), keyData); err != nil {
		return nil, err
	}

	switch keyData.Type {
	case KeyRSA:
		return UnmarshalRSAKey(keyData)

	default:
		return nil, fmt.Errorf("Unknown key type %s", keyData.Type)
	}
}
