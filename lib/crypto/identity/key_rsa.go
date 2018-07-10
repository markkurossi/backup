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

var (
	label = []byte("Backup Zone Secret")
)

type rsaPrivateKey struct {
	name    string
	private *rsa.PrivateKey
}

func (key *rsaPrivateKey) Name() string {
	return key.name
}

func (key *rsaPrivateKey) Type() KeyType {
	return KeyRSAPrivateKey
}

func (key *rsaPrivateKey) Size() int {
	return key.private.PublicKey.N.BitLen()
}

func (key *rsaPrivateKey) ID() string {
	return keyID(&key.private.PublicKey)
}

func (key *rsaPrivateKey) Marshal() ([]byte, error) {
	keyData := &KeyData{
		Name: key.name,
		Type: KeyRSAPrivateKey,
		Data: x509.MarshalPKCS1PrivateKey(key.private),
	}
	return encoding.Marshal(keyData)
}

func (key *rsaPrivateKey) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, key.private, ciphertext,
		label)
}

func (key *rsaPrivateKey) PublicKey() PublicKey {
	return &rsaPublicKey{
		name:   key.name,
		public: &key.private.PublicKey,
	}
}

type rsaPublicKey struct {
	name   string
	public *rsa.PublicKey
}

func (key *rsaPublicKey) Name() string {
	return key.name
}

func (key *rsaPublicKey) Type() KeyType {
	return KeyRSAPublicKey
}

func (key *rsaPublicKey) Size() int {
	return key.public.N.BitLen()
}

func (key *rsaPublicKey) ID() string {
	return keyID(key.public)
}

func (key *rsaPublicKey) Marshal() ([]byte, error) {
	keyData := &KeyData{
		Name: key.name,
		Type: KeyRSAPublicKey,
		Data: x509.MarshalPKCS1PublicKey(key.public),
	}
	return encoding.Marshal(keyData)
}

func (key *rsaPublicKey) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key.public, msg, label)
}

func keyID(key *rsa.PublicKey) string {
	data := x509.MarshalPKCS1PublicKey(key)
	sum := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%s", base64.URLEncoding.EncodeToString(sum[:]))
}

func NewRSAKey(name string, bits int) (PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return &rsaPrivateKey{
		name:    name,
		private: key,
	}, nil
}

func UnmarshalRSAPrivateKey(data *KeyData) (PrivateKey, error) {
	key, err := x509.ParsePKCS1PrivateKey(data.Data)
	if err != nil {
		return nil, err
	}

	return &rsaPrivateKey{
		name:    data.Name,
		private: key,
	}, nil
}

func UnmarshalRSAPublicKey(data *KeyData) (PublicKey, error) {
	key, err := x509.ParsePKCS1PublicKey(data.Data)
	if err != nil {
		return nil, err
	}

	return &rsaPublicKey{
		name:   data.Name,
		public: key,
	}, nil
}
