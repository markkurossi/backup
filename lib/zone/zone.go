//
// zone.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package zone

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/markkurossi/backup/lib/crypto/identity"
	"github.com/markkurossi/backup/lib/local"
)

var zoneDirs = []string{
	"identities",
	"objects",
}

type Zone struct {
	secret []byte
	local  *local.Root
	Name   string
}

func (zone *Zone) identities() string {
	return fmt.Sprintf("%s/identities", zone.Name)
}

func (zone *Zone) objects() string {
	return fmt.Sprintf("%s/objects", zone.Name)
}

func (zone *Zone) AddIdentity(key identity.PublicKey) error {
	encrypted, err := key.Encrypt(zone.secret)
	if err != nil {
		return err
	}
	return zone.local.Add(zone.identities(), key.ID(), encrypted)
}

func Create(local *local.Root, name string) (*Zone, error) {
	secret := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		return nil, err
	}

	for _, dir := range zoneDirs {
		err := local.Mkdir(fmt.Sprintf("%s/%s", name, dir))
		if err != nil {
			return nil, err
		}
	}

	return &Zone{
		secret: secret,
		local:  local,
		Name:   name,
	}, nil
}

func Open(local *local.Root, name string, keys []identity.PrivateKey) (
	*Zone, error) {

	zone := &Zone{
		local: local,
		Name:  name,
	}

	// Get zone identities.
	identities, err := local.GetAll(zone.identities())
	if err != nil {
		return nil, err
	}

	// Do we have an identity to open the zone?
	for _, key := range keys {
		data, ok := identities[key.ID()]
		if !ok {
			continue
		}
		secret, err := key.Decrypt(data)
		if err != nil {
			continue
		}
		zone.secret = secret

		return zone, nil
	}
	return nil, fmt.Errorf("No key to open zone '%s'", name)
}
