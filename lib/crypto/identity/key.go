//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package identity

import (
	"bytes"
	"fmt"

	"github.com/markkurossi/backup/lib/encoding"
)

// Key implements a keypair.
type Key interface {
	Name() string
	Type() KeyType
	Size() int
	ID() string
	Marshal() ([]byte, error)
}

// PrivateKey implements a private key.
type PrivateKey interface {
	Key
	Decrypt(ciphertext []byte) ([]byte, error)
	PublicKey() PublicKey
}

// PublicKey implements a public key.
type PublicKey interface {
	Key
	Encrypt(msg []byte) ([]byte, error)
}

// KeyType defines a key type.
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

// Supported key types.
const (
	KeyRSAPrivateKey KeyType = iota
	KeyRSAPublicKey
)

// KeyData implements a keypair.
type KeyData struct {
	Name string
	Type KeyType
	Data []byte
}

// Unmarshal decodes key from the data.
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
		return nil, fmt.Errorf("invalid key type %s", keyData.Type)
	}
}

// UnmarshalPrivateKey decodes private key from the data.
func UnmarshalPrivateKey(data []byte) (PrivateKey, error) {
	key, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	private, ok := key.(PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid private key: %T", key)
	}
	return private, nil
}

// UnmarshalPublicKey decodes public key from the data.
func UnmarshalPublicKey(data []byte) (PublicKey, error) {
	key, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	private, ok := key.(PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key: %T", key)
	}
	return private, nil
}
