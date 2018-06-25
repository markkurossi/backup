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
	"github.com/markkurossi/backup/tree"
)

func main() {
	fileRoot := flag.String("file-root", "", "Filesystem storage directory")
	list := flag.String("list", "", "List tree ID")
	flag.Parse()

	var storageWriter storage.Writer
	var storageReader storage.Reader
	var err error

	if len(*fileRoot) > 0 {
		var file *storage.File
		file, err = storage.NewFile(*fileRoot)
		storageWriter = file
		storageReader = file
	} else {
		null := storage.NewNull()
		storageWriter = null
		storageReader = null
	}
	if err != nil {
		fmt.Printf("Failed to initialize storage: %s\n", err)
		return
	}

	if len(*list) > 0 {
		id, err := tree.IDFromString(*list)
		if err != nil {
			fmt.Printf("Invalid tree ID '%s': %s\n", *list, err)
			return
		}
		err = remote.List(id, storageReader)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else {
		for _, file := range flag.Args() {
			id, err := local.Traverse(file, storageWriter)
			if err != nil {
				log.Printf("%s\n", err)
			}
			if id != nil {
				fmt.Printf("Tree ID: %s\n", id.ToFullString())
			}
		}
	}
}
