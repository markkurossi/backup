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

	"github.com/markkurossi/backup/lib/objtree"
)

func cmdLs() {
	snapshot := flag.Bool("s", false, "List snapshots.")
	long := flag.Bool("l", false, "List in long format.")
	debug := flag.Bool("d", false, "Enable debugging.")
	flag.Parse()

	if *debug {
		fmt.Printf("Debugging enabled\n")
	}

	z, _ := openZone("default")
	fmt.Printf("Zone '%s' opened\n", z.Name)

	var err error

	id := z.HeadID

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
