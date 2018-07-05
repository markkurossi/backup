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
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/markkurossi/backup/lib/encoding"
)

type rsaKey struct {
	name    string
	private *rsa.PrivateKey
}

func (key *rsaKey) Name() string {
	return key.name
}

func (key *rsaKey) ID() string {
	data := x509.MarshalPKCS1PublicKey(&key.private.PublicKey)
	sum := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%s", base64.StdEncoding.EncodeToString(sum[:]))
}

func (key *rsaKey) Marshal() ([]byte, error) {
	keyData := &KeyData{
		Name: key.name,
		Type: KeyRSA,
		Data: x509.MarshalPKCS1PrivateKey(key.private),
	}
	return encoding.Marshal(keyData)
}

func NewRSAKey(name string, bits int) (Key, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return &rsaKey{
		name:    name,
		private: key,
	}, nil
}

func UnmarshalRSAKey(data *KeyData) (Key, error) {
	key, err := x509.ParsePKCS1PrivateKey(data.Data)
	if err != nil {
		return nil, err
	}

	return &rsaKey{
		name:    data.Name,
		private: key,
	}, nil
}
