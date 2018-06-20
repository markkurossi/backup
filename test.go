//
// test.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/markkurossi/backup/local"
	"github.com/markkurossi/backup/storage"
)

func main() {
	flag.Parse()

	if false {
		hash := sha256.New()
		buf := make([]byte, 4096)

		for _, file := range flag.Args() {
			f, err := os.Open(file)
			if err != nil {
				log.Printf("Failed to open file '%s': %s\n", file, err)
				continue
			}
			for {
				read, err := f.Read(buf)
				if read == 0 {
					if err == io.EOF {
						fmt.Printf("%x  %s\n", hash.Sum(nil), file)
						break
					}
					log.Printf("Read failed: %s\n", err)
					break
				} else {
					hash.Write(buf[:read])
				}
			}

			f.Close()
		}
	} else {
		for _, file := range flag.Args() {
			id, err := local.Traverse(file, storage.NewNull())
			if err != nil {
				log.Printf("%s\n", err)
			}
			fmt.Printf("Tree ID: %s\n", id)
		}
	}
}
