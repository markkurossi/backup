//
// cmd_zone.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"

	"github.com/markkurossi/backup/lib/crypto/identity"
)

func cmdZone() {
	addID := flag.String("a", "", "Add identity")
	flag.Parse()

	z, _ := openZone("default")
	fmt.Printf("Zone '%s' opened\n", z.Name)

	if len(*addID) > 0 {
		var key identity.PublicKey

		switch *addID {
		case "null":
			k, err := identity.GetNull()
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			key = k.PublicKey()
		}

		if key != nil {
			z.AddIdentity(key)
		}
	}
}
