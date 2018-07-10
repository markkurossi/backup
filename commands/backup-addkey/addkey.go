//
// addkey.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/markkurossi/backup/lib/agent"
	"github.com/markkurossi/backup/lib/crypto/identity"
	"github.com/markkurossi/backup/lib/util"
)

const (
	sockEnv = "BACKUP_AGENT_SOCK"
)

func main() {
	address := flag.String("a", "", "Agent UNIX-domain socket address.")
	addAll := flag.Bool("A", false,
		"Add all identities from your identity storage.")
	list := flag.Bool("l", false, "List all keys store in key agent.")
	flag.Parse()

	var path string

	if len(*address) == 0 {
		var ok bool
		path, ok = os.LookupEnv(sockEnv)
		if !ok {
			log.Fatalf("Agent socket environment variable %s not set\n",
				sockEnv)
		}
	} else {
		path = *address
	}

	client, err := agent.NewClient(path)
	if err != nil {
		fmt.Printf("Failed to connect to agent '%s': %s\n", path, err)
		return
	}

	if *list {
		keys, err := client.ListKeys()
		if err != nil {
			fmt.Printf("Failed to list keys: %s\n", err)
			return
		}
		for _, k := range keys {
			fmt.Printf("%s-%d %s %s\n",
				typeName(k.Type()), k.Size(), k.ID(), k.Name())
			msg := []byte("Hello, world!")
			pub := k.PublicKey()
			cipher, err := pub.Encrypt(msg)
			if err != nil {
				log.Fatalf("Failed to encrypt with public key: %s\n", err)
			}
			plain, err := k.Decrypt(cipher)
			if err != nil {
				log.Fatalf("Failed to decrypt with private key: %s\n", err)
			}
			if !bytes.Equal(msg, plain) {
				log.Fatalf("Data mismatch\n")
			}
		}
	}

	if *addAll {
		user, err := user.Current()
		if err != nil {
			log.Fatalf("Failed to get current user: %s\n", err)
		}
		storage := identity.NewStorage(user)
		if err := storage.Open(); err != nil {
			log.Fatalf("Failed to open identity storage %s: %s\n",
				storage.Dir, err)
		}
		keys, err := storage.List()
		if err != nil {
			log.Fatalf("Failed to list keys: %s\n", err)
		}
		for _, keyInfo := range keys {
			passphrase, err := util.ReadPassphrase(
				fmt.Sprintf("Enter passphrase for key '%s'", keyInfo.Name),
				false)
			if err != nil {
				log.Fatalf("%s\n", err)
			}
			key, err := storage.Load(keyInfo.ID, passphrase)
			if err != nil {
				fmt.Printf("Failed to load key: %s\n", err)
				continue
			}
			err = client.AddKey(key)
			if err != nil {
				fmt.Printf("Failed to add key: %s\n", err)
			}
		}
	}

	for _, cmd := range flag.Args() {
		fmt.Printf("Command: %s\n", cmd)
	}
}

func typeName(keyType identity.KeyType) string {
	switch keyType {
	case identity.KeyRSAPrivateKey, identity.KeyRSAPublicKey:
		return "RSA"

	default:
		return keyType.String()
	}
}
