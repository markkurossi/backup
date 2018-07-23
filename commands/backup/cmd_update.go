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
	"time"

	"github.com/markkurossi/backup/lib/local"
	"github.com/markkurossi/backup/lib/tree"
)

func cmdUpdate() {
	debug := flag.Bool("d", false, "Enable debugging.")
	flag.Parse()

	if *debug {
		fmt.Printf("Debugging enabled\n")
	}

	z := openZone("default")
	fmt.Printf("Zone '%s' opened\n", z.Name)

	id, err := local.Traverse(z.Local.Root, z)
	if err != nil {
		fmt.Printf("Failed to traverse directory '%s': %s\n", z.Local.Root, err)
		os.Exit(1)
	}
	if id.Undefined() {
		fmt.Printf("Zone root '%s' is not a directory\n", z.Local.Root)
		os.Exit(1)
	}
	fmt.Printf("Tree ID: %s\n", id)
	if z.Written > 0 {
		fmt.Printf("Data size: %d, saved %d (%.0f%%)\n", z.Written, z.Saved,
			float64(z.Saved)/float64(z.Written)*100.0)
	}

	if z.Head != nil && id.Equal(z.Head.Root) {
		fmt.Printf("No changes\n")
		os.Exit(0)
	}

	snapshot := tree.NewSnapshot()
	snapshot.Timestamp = time.Now().UnixNano()
	snapshot.Size = tree.FileSize(z.Written)
	snapshot.Root = id
	if z.Head != nil {
		snapshot.Parent = z.HeadID
	}

	data, err := snapshot.Serialize()
	if err != nil {
		fmt.Printf("Failed to serialize snapshot: %s\n", err)
		os.Exit(1)
	}
	headID, err := z.Write(data)
	if err != nil {
		fmt.Printf("Failed to write snapshot: %s\n", err)
		os.Exit(1)
	}

	z.Head = snapshot
	z.HeadID = headID

	err = z.SetRootPointer(headID)
	if err != nil {
		fmt.Printf("Failed to save snapshot: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Snapshot: %s\n", headID)
}
