//
// storage.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//
// Identity storage.
//

package identity

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

func NewStorage(user *user.User) *Storage {
	return &Storage{
		Dir: fmt.Sprintf("%s/.backup/identities", user.HomeDir),
	}
}

type Storage struct {
	Dir string
}

func (s *Storage) Open() error {
	return os.MkdirAll(s.Dir, 0700)
}

func (s *Storage) Save(key Key, passphrase []byte) error {
	data, err := key.Marshal()
	if err != nil {
		return err
	}
	encrypted, err := Encrypt(data, EncrAES128GCM, passphrase,
		KDFPBKDF24096SHA256)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%s", s.Dir, key.ID())
	return ioutil.WriteFile(path, encrypted, 0700)
}
