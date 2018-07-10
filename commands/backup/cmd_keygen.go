//
// cmd_keygen.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/markkurossi/backup/lib/crypto/identity"
	"github.com/markkurossi/backup/lib/util"
)

func cmdKeygen() {
	flag.Parse()

	user, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get current user: %s\n", err)
		os.Exit(1)
	}
	bits := 4096
	fmt.Printf("Creating %d bit RSA key...\n", bits)
	key, err := identity.NewRSAKey(user.Username, bits)
	if err != nil {
		fmt.Printf("Identity key generation failed: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created identity key %s\n", key.ID())

	storage := identity.NewStorage(user)
	if err := storage.Open(); err != nil {
		fmt.Printf("Failed to open identity storage %s: %s\n",
			storage.Dir, err)
		os.Exit(1)
	}

	passphrase, err := util.ReadPassphrase("Enter passphrase for the key", true)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	err = storage.Save(key, passphrase)
	if err != nil {
		fmt.Printf("Failed to save key: %s\n", err)
		os.Exit(1)
	}
}
