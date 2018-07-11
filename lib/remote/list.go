//
// list.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package remote

import (
	"fmt"

	"github.com/markkurossi/backup/lib/storage"
	"github.com/markkurossi/backup/lib/tree"
)

func List(root *storage.ID, st storage.Accessor) error {
	return list("", true, root, st)
}

func nest(indent string, isLast bool) string {
	if isLast {
		return indent + "    "
	} else {
		return indent + "|   "
	}
}

func list(indent string, verbose bool, root *storage.ID,
	st storage.Accessor) error {
	element, err := tree.Deserialize(root, st)
	if err != nil {
		return err
	}

	if element.IsDir() {
		count := len(element.Directory().Entries)
		for idx, e := range element.Directory().Entries {
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
				fmt.Printf("\t%s\t%s", e.Mode, e.Entry)
			}
			fmt.Println()
			err := list(nest(indent, isLast), verbose, e.Entry, st)
			if err != nil {
				return err
			}
		}
	}

	return nil
}