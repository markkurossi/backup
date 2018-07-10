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
	Name() string
	Type() KeyType
	Size() int
	ID() string
	Marshal() ([]byte, error)
}

type PrivateKey interface {
	Key
	Decrypt(ciphertext []byte) ([]byte, error)
	PublicKey() PublicKey
}

type PublicKey interface {
	Key
	Encrypt(msg []byte) ([]byte, error)
}

type KeyType int

func (t KeyType) String() string {
	switch t {
	case KeyRSAPrivateKey:
		return "rsa-private-key"

	case KeyRSAPublicKey:
		return "rsa-public-key"

	default:
		return fmt.Sprintf("{KeyType %d}", t)
	}
}

const (
	KeyRSAPrivateKey KeyType = 0
	KeyRSAPublicKey          = 1
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
	case KeyRSAPrivateKey:
		return UnmarshalRSAPrivateKey(keyData)

	case KeyRSAPublicKey:
		return UnmarshalRSAPublicKey(keyData)

	default:
		return nil, fmt.Errorf("Invalid key type %s", keyData.Type)
	}
}

func UnmarshalPrivateKey(data []byte) (PrivateKey, error) {
	key, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	private, ok := key.(PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Invalid private key: %T", key)
	}
	return private, nil
}

func UnmarshalPublicKey(data []byte) (PublicKey, error) {
	key, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	private, ok := key.(PublicKey)
	if !ok {
		return nil, fmt.Errorf("Invalid public key: %T", key)
	}
	return private, nil
}
