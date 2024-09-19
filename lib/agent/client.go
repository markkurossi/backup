//
// Copyright (c) 2018-2024 Markku Rossi
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

// Client implements an agent connection.
type Client struct {
	conn net.Conn
}

// NewClient creates a new client and connects to the agent.
func NewClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// AddKey adds the identity key to the agent.
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
		return fmt.Errorf("unsupported agent message '%s'", msg.Type())
	}
}

// ListKeys lists all identity keys in the agent.
func (c *Client) ListKeys() ([]identity.PrivateKey, error) {
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
		var keys []identity.PrivateKey
		for _, ki := range m.Keys {
			pub, err := identity.UnmarshalPublicKey(ki.PublicKey)
			if err != nil {
				return nil, err
			}

			keys = append(keys, &proxyKey{
				client:    c,
				info:      ki,
				publicKey: pub,
			})
		}
		return keys, nil

	default:
		return nil, fmt.Errorf("unsupported agent message '%s'", msg.Type())
	}
}

type proxyKey struct {
	client    *Client
	info      KeyInfo
	publicKey identity.PublicKey
}

func (key *proxyKey) Name() string {
	return key.info.Name
}

func (key *proxyKey) Type() identity.KeyType {
	return key.info.Type
}

func (key *proxyKey) Size() int {
	return key.info.Size
}

func (key *proxyKey) ID() string {
	return key.info.ID
}

func (key *proxyKey) Marshal() ([]byte, error) {
	return nil, errors.New("Marshal not implemented yet")
}

func (key *proxyKey) Decrypt(ciphertext []byte) ([]byte, error) {
	msg, err := RPC(key.client.conn, &MsgDecrypt{
		MsgHdr: MsgHdr{
			t: Decrypt,
		},
		KeyID: key.info.ID,
		Data:  ciphertext,
	})
	if err != nil {
		return nil, err
	}
	switch m := msg.(type) {
	case *MsgError:
		return nil, errors.New(m.Message)

	case *MsgDecrypted:
		return m.Data, nil

	default:
		return nil, fmt.Errorf("unsupported agent message '%s'", msg.Type())
	}
}

func (key *proxyKey) PublicKey() identity.PublicKey {
	return key.publicKey
}
