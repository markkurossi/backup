//
// key.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markkurossi/backup/lib/agent"
)

const (
	sockEnv = "BACKUP_AGENT_SOCK"
)

func main() {
	path, ok := os.LookupEnv(sockEnv)
	if !ok {
		fmt.Printf("Agent socket environment variable %s not set\n", sockEnv)
	}

	client, err := agent.NewClient(path)
	if err != nil {
		fmt.Printf("Failed to connect to agent '%s': %s\n", path, err)
		return
	}

	keys, err := client.ListKeys()
	if err != nil {
		fmt.Printf("Failed to list keys: %s\n", err)
		return
	}
	log.Printf("Keys: %v\n", keys)
}
