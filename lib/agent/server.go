//
// server.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package agent

import (
	"io"
	"net"
)

type Server struct {
	listener net.Listener
}

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

func NewServer(listener net.Listener) *Server {
	return &Server{
		listener: listener,
	}
}

type Connection struct {
	conn net.Conn
	C    chan Msg
}

func (c *Connection) SendOK() error {
	return SendMessage(c.conn, &MsgOK{
		MsgHdr: MsgHdr{
			t: OK,
		},
	})
}

func (c *Connection) SendError(msg string) error {
	return SendMessage(c.conn, &MsgError{
		MsgHdr: MsgHdr{
			t: Error,
		},
		Message: msg,
	})
}

func (c *Connection) SendKeys(keys []KeyInfo) error {
	return SendMessage(c.conn, &MsgKeys{
		MsgHdr: MsgHdr{
			t: Keys,
		},
		Keys: keys,
	})
}

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
