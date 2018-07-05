//
// key_rsa.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package identity

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/markkurossi/backup/lib/encoding"
)

type rsaKey struct {
	name string
	key  *rsa.PrivateKey
}

func (key *rsaKey) Marshal() ([]byte, error) {
	keyData := &KeyData{
		Name: key.name,
		Type: KeyRSA,
		Data: x509.MarshalPKCS1PrivateKey(key.key),
	}
	return encoding.Marshal(keyData)
}

func NewRSAKey(name string, bits int) (Key, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return &rsaKey{
		name: name,
		key:  key,
	}, nil
}

func UnmarshalRSAKey(data *KeyData) (Key, error) {
	key, err := x509.ParsePKCS1PrivateKey(data.Data)
	if err != nil {
		return nil, err
	}

	return &rsaKey{
		name: data.Name,
		key:  key,
	}, nil
}
