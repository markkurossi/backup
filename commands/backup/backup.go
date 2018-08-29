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
	"github.com/markkurossi/backup/lib/crypto/zone"
	"github.com/markkurossi/backup/lib/persistence"
)

var commands = map[string]func(){
	"add-key": cmdAddKey,
	"init":    cmdInit,
	"keygen":  cmdKeygen,
	"ls":      cmdLs,
	"update":  cmdUpdate,
	"zone":    cmdZone,
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

func openZone(name string) (*zone.Zone, string) {
	connectAgent()

	keys, err := client.ListKeys()
	if err != nil {
		fmt.Printf("Failed to get identity keys: %s\n", err)
		os.Exit(1)
	}
	if len(keys) == 0 {
		fmt.Printf("No identity keys defined\n")
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current working directory: %s\n", err)
		os.Exit(1)
	}
	root, err := persistence.OpenFilesystem(fmt.Sprintf("%s/.backup", wd))
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	z, err := zone.Open(root, "default", keys)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	return z, wd
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
		os.Exit(1)
	}
	flag.CommandLine = flag.NewFlagSet(fmt.Sprintf("backup %s", os.Args[0]),
		flag.ExitOnError)
	fn()
}
