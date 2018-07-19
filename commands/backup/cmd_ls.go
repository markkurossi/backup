//
// cmd_ls.go
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

	"github.com/markkurossi/backup/lib/objtree"
)

func cmdLs() {
	debug := flag.Bool("d", false, "Enable debugging.")
	root := flag.String("r", "", "Tree root ID")
	flag.Parse()

	if *debug {
		fmt.Printf("Debugging enabled\n")
	}

	if len(*root) == 0 {
		fmt.Printf("No tree root ID defined\n")
		os.Exit(1)
	}

	z := openZone("default")
	fmt.Printf("Zone '%s' opened\n", z.Name)

	id, err := z.ResolveID(*root)
	if err != nil {
		fmt.Printf("Invalid tree ID '%s': %s\n", *root, err)
		os.Exit(1)
	}

	err = objtree.List(id, z)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
