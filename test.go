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
	"github.com/markkurossi/backup/storage"
)

func main() {
	fileRoot := flag.String("file-root", "", "Filesystem storage directory")
	flag.Parse()

	var storageWriter storage.Writer
	var err error

	if fileRoot != nil {
		storageWriter, err = storage.NewFile(*fileRoot)
	} else {
		storageWriter = storage.NewNull()
	}
	if err != nil {
		fmt.Printf("Failed to initialize storage: %s\n", err)
		return
	}

	for _, file := range flag.Args() {
		id, err := local.Traverse(file, storageWriter)
		if err != nil {
			log.Printf("%s\n", err)
		}
		fmt.Printf("Tree ID: %s\n", id)
	}
}
