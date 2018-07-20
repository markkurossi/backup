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
	"github.com/markkurossi/backup/lib/storage"
)

func cmdLs() {
	snapshot := flag.Bool("s", false, "List snapshots.")
	long := flag.Bool("l", false, "List in long format.")
	debug := flag.Bool("d", false, "Enable debugging.")
	root := flag.String("r", "", "Root object ID.")
	flag.Parse()

	if *debug {
		fmt.Printf("Debugging enabled\n")
	}

	z := openZone("default")
	fmt.Printf("Zone '%s' opened\n", z.Name)

	var id storage.ID
	var err error

	if len(*root) > 0 {
		id, err = z.ResolveID(*root)
		if err != nil {
			fmt.Printf("Invalid tree ID '%s': %s\n", *root, err)
			os.Exit(1)
		}
	} else {
		id = z.HeadID
	}

	if *snapshot {
		// List snapshots.
		err = objtree.ListSnapshots(id, z, *long)
	} else {
		// List files.
		err = objtree.List(id, z, *long)
	}
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
