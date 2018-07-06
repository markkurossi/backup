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
	"io/ioutil"
	"log"
	"os"
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

	keyData, err := key.Marshal()
	if err != nil {
		log.Fatalf("Failed to marshal key: %s\n", err)
	}

	keyDir := fmt.Sprintf("%s/.backup/identities", user.HomeDir)
	err = os.MkdirAll(keyDir, 0700)
	if err != nil {
		log.Fatalf("Failed to create identity directory %s: %s\n", keyDir, err)
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

	// Encrypt key.
	encrypted, err := identity.Encrypt(keyData, identity.EncrAES128GCM,
		passphrase, identity.KDFPBKDF24096SHA256)
	if err != nil {
		log.Fatalf("Failed to encrypt key: %s\n", err)
	}

	// And save it.
	keyPath := fmt.Sprintf("%s/%s", keyDir, key.ID())
	err = ioutil.WriteFile(keyPath, encrypted, 0700)
	if err != nil {
		log.Fatalf("Failed to save key: %s\n", err)
	}
}
