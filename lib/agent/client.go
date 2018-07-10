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

	"github.com/markkurossi/backup/lib/crypto/identity"
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

func (c *Client) AddKey(key identity.Key) error {
	data, err := key.Marshal()
	if err != nil {
		return err
	}
	msg, err := RPC(c.conn, &MsgAddKey{
		MsgHdr: MsgHdr{
			t: AddKey,
		},
		Data: data,
	})
	if err != nil {
		return err
	}
	switch m := msg.(type) {
	case *MsgError:
		return errors.New(m.Message)

	case *MsgOK:
		return nil

	default:
		return fmt.Errorf("Unsupported agent message '%s'", msg.Type())
	}
}

func (c *Client) ListKeys() ([]KeyInfo, error) {
	msg, err := RPC(c.conn, &MsgListKeys{
		MsgHdr: MsgHdr{
			t: ListKeys,
		},
	})
	if err != nil {
		return nil, err
	}
	switch m := msg.(type) {
	case *MsgError:
		return nil, errors.New(m.Message)

	case *MsgKeys:
		return m.Keys, nil

	default:
		return nil, fmt.Errorf("Unsupported agent message '%s'", msg.Type())
	}
}
