//
// backup.go
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
)

var commands = map[string]func(){
	"init":    cmdInit,
	"keygen":  cmdKeygen,
	"add-key": cmdAddKey,
}

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output.")
	flag.Parse()

	if *verbose {
		fmt.Printf("Verbose mode enabled\n")
	}
	if len(flag.Args()) == 0 {
		flag.Usage()
		fmt.Printf("Possible commands are:\n")
		for key := range commands {
			fmt.Printf(" - %s\n", key)
		}
		return
	}
	os.Args = flag.Args()
	fn, ok := commands[flag.Arg(0)]
	if !ok {
		fmt.Printf("Unknown command: %s\n", flag.Arg(0))
	}
	flag.CommandLine = flag.NewFlagSet(fmt.Sprintf("backup %s", os.Args[0]),
		flag.ExitOnError)
	fn()
}
