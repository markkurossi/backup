//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package agent

import (
	"io"
	"net"
)

// Server implements an agent server.
type Server struct {
	listener net.Listener
}

// Accept accepts a new client connection.
func (s *Server) Accept() (*Connection, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	c := &Connection{
		conn: conn,
		C:    make(chan Msg),
	}
	go c.messageLoop()

	return c, nil
}

// NewServer creates a new agent server.
func NewServer(listener net.Listener) *Server {
	return &Server{
		listener: listener,
	}
}

// Connection implements a client connection.
type Connection struct {
	conn net.Conn
	C    chan Msg
}

// SendOK sends to success message to the connection.
func (c *Connection) SendOK() error {
	return SendMessage(c.conn, &MsgOK{
		MsgHdr: MsgHdr{
			t: OK,
		},
	})
}

// SendError sends the error message msg to the connection.
func (c *Connection) SendError(msg string) error {
	return SendMessage(c.conn, &MsgError{
		MsgHdr: MsgHdr{
			t: Error,
		},
		Message: msg,
	})
}

// SendKeys sends the identity keys to the connection.
func (c *Connection) SendKeys(keys []KeyInfo) error {
	return SendMessage(c.conn, &MsgKeys{
		MsgHdr: MsgHdr{
			t: Keys,
		},
		Keys: keys,
	})
}

// SendDecrypted sends decrypted data to the connection.
func (c *Connection) SendDecrypted(data []byte) error {
	return SendMessage(c.conn, &MsgDecrypted{
		MsgHdr: MsgHdr{
			t: Decrypted,
		},
		Data: data,
	})
}

func (c *Connection) messageLoop() {
	for {
		msg, err := ReceiveMessage(c.conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			c.SendError(err.Error())
		} else {
			c.C <- msg
		}
	}

	close(c.C)
}
