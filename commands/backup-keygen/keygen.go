//
// keygen.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"bytes"
	"fmt"
	"log"
	"os/user"
	"syscall"

	"github.com/markkurossi/backup/lib/crypto/identity"
	"golang.org/x/crypto/ssh/terminal"
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

	var passphrase []byte

	for {
		fmt.Print("Enter passphrase for the identity key: ")
		first, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\n%s\n", err)
			return
		}
		fmt.Print("\nEnter same passphrase again: ")
		second, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\n%s\n", err)
			return
		}
		if !bytes.Equal(first, second) {
			fmt.Print("\nPassphrases do not match\n")
			continue
		}
		fmt.Print("\n")
		passphrase = first
		break
	}

	err = storage.Save(key, passphrase)
	if err != nil {
		log.Fatalf("Failed to save key: %s\n", err)
	}
}
