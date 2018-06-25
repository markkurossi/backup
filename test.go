//
// test.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/markkurossi/backup/local"
	"github.com/markkurossi/backup/remote"
	"github.com/markkurossi/backup/storage"
)

func main() {
	fileRoot := flag.String("file-root", "", "Filesystem storage directory")
	list := flag.String("list", "", "List tree ID")
	flag.Parse()

	var storageAccessor storage.Accessor
	var err error

	if len(*fileRoot) > 0 {
		storageAccessor, err = storage.NewFile(*fileRoot)
	} else {
		storageAccessor = storage.NewNull()
	}
	if err != nil {
		fmt.Printf("Failed to initialize storage: %s\n", err)
		return
	}

	if len(*list) > 0 {
		id, err := storage.IDFromString(*list)
		if err != nil {
			fmt.Printf("Invalid tree ID '%s': %s\n", *list, err)
			return
		}
		err = remote.List(id, storageAccessor)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else {
		for _, file := range flag.Args() {
			id, err := local.Traverse(file, storageAccessor)
			if err != nil {
				log.Printf("%s\n", err)
			}
			if id != nil {
				fmt.Printf("Tree ID: %s\n", id.ToFullString())
			}
		}
	}
}
