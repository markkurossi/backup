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
	"github.com/markkurossi/backup/lib/crypto/identity"
)

var (
	identities = make(map[string]identity.PrivateKey)
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
messageLoop:
	for msg := range c.C {
		log.Printf("Message: %v\n", msg)
		switch m := msg.(type) {
		case *agent.MsgAddKey:
			key, err := identity.UnmarshalPrivateKey(m.Data)
			if err != nil {
				txt := fmt.Sprintf("Invalid key data: %s", err)
				log.Printf("%s\n", txt)
				c.SendError(txt)
			} else {
				// Do we already have this key?
				_, ok := identities[key.ID()]
				if !ok {
					identities[key.ID()] = key
				}
				c.SendOK()
			}

		case *agent.MsgListKeys:
			var keys []agent.KeyInfo
			for _, key := range identities {
				pub, err := key.PublicKey().Marshal()
				if err != nil {
					txt := fmt.Sprintf("Failed to marshal public key: %s", err)
					log.Printf("%s\n", err)
					c.SendError(txt)
					continue messageLoop
				}
				keys = append(keys, agent.KeyInfo{
					Name:      key.Name(),
					Type:      key.Type(),
					Size:      key.Size(),
					ID:        key.ID(),
					PublicKey: pub,
				})
			}
			c.SendKeys(keys)

		case *agent.MsgDecrypt:
			key, ok := identities[m.KeyID]
			if !ok {
				txt := fmt.Sprintf("Unknown key: %s", m.KeyID)
				log.Printf("%s\n", txt)
				c.SendError(txt)
				continue
			}
			data, err := key.Decrypt(m.Data)
			if err != nil {
				c.SendError(err.Error())
				continue
			}
			c.SendDecrypted(data)

		default:
			txt := fmt.Sprintf("Unsupported client message '%s'", msg.Type())
			log.Printf("%s\n", txt)
			c.SendError(txt)
		}
	}
}
