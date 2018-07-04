//
// key-agent.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/markkurossi/backup/lib/agent"
)

func main() {
	bindAddress := flag.String("a", "",
		"Bind address for the agent UNIX-domain socket.")
	flag.Parse()

	if len(*bindAddress) == 0 {
		dir, err := ioutil.TempDir("", "backup-agent")
		if err != nil {
			log.Fatalf("Failed to create temporary directory: %s\n", err)
		}
		defer os.RemoveAll(dir)
		*bindAddress = fmt.Sprintf("%s/agent.%d", dir, os.Getpid())
	}
	defer os.Remove(*bindAddress)

	os.Remove(*bindAddress)
	umask := syscall.Umask(0077)
	listener, err := net.Listen("unix", *bindAddress)
	syscall.Umask(umask)
	if err != nil {
		log.Fatalf("Failed to create listener %s: %s\n", *bindAddress, err)
	}
	fmt.Printf("export BACKUP_AGENT_SOCK=%s\n", *bindAddress)

	server := agent.NewServer(listener)
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %s\n", err)
			break
		}
		go handleConnection(conn)
	}
}

func handleConnection(c *agent.Connection) {
	for msg := range c.C {
		log.Printf("Message: %v\n", msg)
		switch msg.(type) {
		default:
			txt := fmt.Sprintf("Unsupported client message '%s'", msg.Type())
			log.Printf("%s\n", txt)
			c.SendError(txt)
		}
	}
}
