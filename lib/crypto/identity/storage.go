//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//
// Identity storage.
//

package identity

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/markkurossi/backup/lib/encoding"
)

// NewStorage creates a new storage for the user.
func NewStorage(user *user.User) *Storage {
	return &Storage{
		Dir: fmt.Sprintf("%s/.backup.d/identities", user.HomeDir),
	}
}

// Storage implements an identity storage.
type Storage struct {
	Dir string
}

// Open opens the storage.
func (s *Storage) Open() error {
	return os.MkdirAll(s.Dir, 0700)
}

// KeyInfo provides information about an identity key.
type KeyInfo struct {
	ID   string
	Name string
}

// List lists all identity keys.
func (s *Storage) List() ([]KeyInfo, error) {
	info, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}
	var keys []KeyInfo
	for _, fi := range info {
		if fi.IsDir() {
			continue
		}
		data, err := s.loadKeyData(fi.Name())
		if err != nil {
			log.Printf("Failed to read identity key %s: %s\n", fi.Name(), err)
			continue
		}
		enc := new(EncryptedKey)
		if err = encoding.Unmarshal(bytes.NewReader(data), enc); err != nil {
			log.Printf("Skipping unparseable identity key %s (%s)\n",
				fi.Name(), err)
			continue
		}

		keys = append(keys, KeyInfo{
			ID:   fi.Name(),
			Name: enc.Name,
		})
	}
	return keys, nil
}

// Load loads the key id that is encrypted with the passphrase.
func (s *Storage) Load(id string, passphrase []byte) (Key, error) {
	encrypted, err := s.loadKeyData(id)
	if err != nil {
		return nil, err
	}
	data, err := Decrypt(encrypted, passphrase)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}

func (s *Storage) loadKeyData(id string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("%s/%s", s.Dir, id))
}

// Save saves they key encrypted with the passphrase.
func (s *Storage) Save(key Key, passphrase []byte) error {
	data, err := key.Marshal()
	if err != nil {
		return err
	}
	encrypted, err := Encrypt(data, EncrAES128GCM, key.Name(), passphrase,
		KDFPBKDF24096SHA256)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%s", s.Dir, key.ID())
	return ioutil.WriteFile(path, encrypted, 0700)
}
