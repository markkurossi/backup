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

	"github.com/markkurossi/backup/lib/agent"
)

var commands = map[string]func(){
	"init":    cmdInit,
	"update":  cmdUpdate,
	"keygen":  cmdKeygen,
	"add-key": cmdAddKey,
}

var address = flag.String("a", "", "Agent UNIX-domain socket address.")
var verbose = flag.Bool("v", false, "Enable verbose output.")

var client *agent.Client

func connectAgent() {
	var path string
	var err error

	if len(*address) == 0 {
		var ok bool
		path, ok = os.LookupEnv(sockEnv)
		if !ok {
			fmt.Printf("Agent socket environment variable %s not set\n",
				sockEnv)
			os.Exit(1)
		}
	} else {
		path = *address
	}

	client, err = agent.NewClient(path)
	if err != nil {
		fmt.Printf("Failed to connect to agent '%s': %s\n", path, err)
		os.Exit(1)
	}
}

func main() {
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
