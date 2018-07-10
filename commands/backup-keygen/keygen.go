//
// keygen.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/markkurossi/backup/lib/crypto/identity"
	"github.com/markkurossi/backup/lib/util"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %s\n", err)
	}
	bits := 4096
	fmt.Printf("Creating %d bit RSA key...\n", bits)
	key, err := identity.NewRSAKey(user.Username, bits)
	if err != nil {
		log.Fatalf("Identity key generation failed: %s\n", err)
	}
	fmt.Printf("Created identity key %s\n", key.ID())

	storage := identity.NewStorage(user)
	if err := storage.Open(); err != nil {
		log.Fatalf("Failed to open identity storage %s: %s\n",
			storage.Dir, err)
	}

	passphrase, err := util.ReadPassphrase("Enter passphrase for the key", true)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	err = storage.Save(key, passphrase)
	if err != nil {
		log.Fatalf("Failed to save key: %s\n", err)
	}
}
