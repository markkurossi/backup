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

type Zone struct {
	secret []byte
	local  *local.Root
	Name   string
}

func (zone *Zone) AddIdentity(key identity.PublicKey) error {
	encrypted, err := key.Encrypt(zone.secret)
	if err != nil {
		return err
	}
	return zone.local.Add(zone.Name, key.ID(), encrypted)
}

func Create(local *local.Root, name string) (*Zone, error) {
	secret := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		return nil, err
	}
	return &Zone{
		secret: secret,
		local:  local,
		Name:   name,
	}, nil
}

func Open(local *local.Root, name string, keys []identity.PrivateKey) (
	*Zone, error) {

	// Get zone identities.
	identities, err := local.GetAll(name)
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

		return &Zone{
			secret: secret,
			local:  local,
			Name:   name,
		}, nil
	}
	return nil, fmt.Errorf("No key to open zone '%s'", name)
}
