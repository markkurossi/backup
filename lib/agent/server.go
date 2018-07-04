//
// server.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package agent

import (
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

func (c *Connection) SendError(msg string) error {
	return SendMessage(c.conn, &MsgError{
		MsgHdr: MsgHdr{
			t: Error,
		},
		Message: msg,
	})
}

func (c *Connection) messageLoop() {
	for {
		msg, err := ReceiveMessage(c.conn)
		if err != nil {
			break
		}

		c.C <- msg
	}

	close(c.C)
}
