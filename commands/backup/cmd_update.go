//
// cmd_update.go
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

	"github.com/markkurossi/backup/lib/local"
	"github.com/markkurossi/backup/lib/zone"
)

func cmdUpdate() {
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

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current working directory: %s\n", err)
		os.Exit(1)
	}
	root, err := local.OpenRoot(wd)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	z, err := zone.Open(root, "default", keys)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Zone '%s' opened\n", z.Name)
}
