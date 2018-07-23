//
// cmd_init.go
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

	"github.com/markkurossi/backup/lib/crypto/zone"
	"github.com/markkurossi/backup/lib/local"
	"github.com/markkurossi/backup/lib/storage"
)

func cmdInit() {
	debug := flag.Bool("d", false, "Enable debugging.")
	flag.Parse()

	if *debug {
		fmt.Printf("Debugging enabled\n")
	}

	connectAgent()

	keys, err := client.ListKeys()
	if err != nil {
		fmt.Printf("Failed to get identity keys: %s\n", err)
		os.Exit(1)
	}
	if len(keys) == 0 {
		fmt.Printf("No identity keys defined\n")
		os.Exit(1)
	}
	// XXX select the default key
	key := keys[0]

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current working directory: %s\n", err)
		os.Exit(1)
	}
	root, err := local.InitRoot(wd)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	z, err := zone.Create(root, "default")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	err = z.AddIdentity(key.PublicKey())
	if err != nil {
		fmt.Printf("Failed to add identity key: %s\n", err)
		os.Exit(1)
	}
	err = z.SetRootPointer(storage.EmptyID)
	if err != nil {
		fmt.Printf("Failed to init root pointer: %s\n", err)
		os.Exit(1)
	}
}
