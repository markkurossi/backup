//
// cmd_init.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
)

func cmdInit() {
	debug := flag.Bool("d", false, "Enable debugging.")
	flag.Parse()

	if *debug {
		fmt.Printf("Debugging enabled\n")
	}
}
