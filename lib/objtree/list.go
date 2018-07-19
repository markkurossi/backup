//
// list.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package objtree

import (
	"fmt"
	"time"

	"github.com/markkurossi/backup/lib/storage"
	"github.com/markkurossi/backup/lib/tree"
)

func List(root storage.ID, st storage.Accessor) error {
	now := time.Now()
	return list(now, "", true, root, st)
}

func nest(indent string, isLast bool) string {
	if isLast {
		return indent + "    "
	} else {
		return indent + "|   "
	}
}

func list(now time.Time, indent string, verbose bool, root storage.ID,
	st storage.Accessor) error {
	element, err := tree.DeserializeID(root, st)
	if err != nil {
		fmt.Printf("Failed to deserialize ID %s: %s\n", root, err)
		return err
	}

	switch el := element.(type) {
	case *tree.Snapshot:
		fmt.Printf("Snapshot %s\n", root)
		fmt.Printf("|-- Created: %s\n", time.Unix(0, el.Timestamp))
		fmt.Printf("|-- Parent : %s\n", el.Parent)
		fmt.Printf("`-- Root   : %s\n", el.Root)
		return list(now, indent+"    ", verbose, el.Root, st)

	case *tree.Directory:
		count := len(el.Entries)
		for idx, e := range el.Entries {
			var in string
			var isLast bool
			if idx+1 == count {
				in = indent + "`-- "
				isLast = true
			} else {
				in = indent + "|-- "
			}
			fmt.Printf("%s%s", in, e.Name)
			if verbose {
				for i := 0; i+len(in)+len(e.Name) < 40; i++ {
					fmt.Printf(" ")
				}
				modified := time.Unix(e.ModTime, 0)
				var modStr string
				if modified.Year() != now.Year() {
					modStr = modified.Format("Jan _2  2006")
				} else {
					modStr = modified.Format("Jan _2 15:04")
				}
				fmt.Printf("\t%s\t%s\t%s", e.Mode, modStr, e.Entry)
			}
			fmt.Println()
			err := list(now, nest(indent, isLast), verbose, e.Entry, st)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
