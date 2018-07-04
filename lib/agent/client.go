//
// client.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package agent

import (
	"errors"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) ListKeys() ([]KeyInfo, error) {
	err := SendMessage(c.conn, &MsgListKeys{
		MsgHdr: MsgHdr{
			t: ListKeys,
		},
	})
	if err != nil {
		return nil, err
	}
	msg, err := ReceiveMessage(c.conn)
	if err != nil {
		return nil, err
	}
	switch m := msg.(type) {
	case *MsgError:
		return nil, errors.New(m.Message)

	default:
		return nil, fmt.Errorf("Unsupported agent message '%s'", msg.Type())
	}

	fmt.Printf("Message: %s\n", msg)
	return nil, nil
}
