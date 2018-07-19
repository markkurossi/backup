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

	"github.com/markkurossi/backup/lib/local"
)

func cmdUpdate() {
	debug := flag.Bool("d", false, "Enable debugging.")
	flag.Parse()

	if *debug {
		fmt.Printf("Debugging enabled\n")
	}

	z := openZone("default")
	fmt.Printf("Zone '%s' opened\n", z.Name)

	id, err := local.Traverse(z.Root(), z)
	if err != nil {
		fmt.Printf("Failed to traverse directory '%s': %s\n", z.Root(), err)
	}
	if id != nil {
		fmt.Printf("Tree ID: %s\n", id.ToFullString())
		if z.Written > 0 {
			fmt.Printf("Data size: %d, saved %d (%.0f%%)\n", z.Written, z.Saved,
				float64(z.Saved)/float64(z.Written)*100.0)
		}
	}
}
