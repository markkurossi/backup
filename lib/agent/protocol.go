//
// protocol.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/markkurossi/backup/lib/crypto/identity"
	"github.com/markkurossi/backup/lib/encoding"
)

type MsgType uint8

const (
	OK       MsgType = 0
	Error            = 1
	AddKey           = 2
	Question         = 3
	Answer           = 4
	ListKeys         = 5
	Keys             = 6
)

var MsgTypeNames = map[MsgType]string{
	OK:       "ok",
	Error:    "error",
	AddKey:   "add-key",
	Question: "question",
	Answer:   "answer",
	ListKeys: "list-keys",
	Keys:     "keys",
}

func (t MsgType) String() string {
	name, ok := MsgTypeNames[t]
	if ok {
		return name
	}
	return fmt.Sprintf("{MsgType %d}", t)
}

type Msg interface {
	SetType(t MsgType)
	Type() MsgType
}

type MsgHdr struct {
	t MsgType `backup:"-"`
}

func (hdr *MsgHdr) SetType(t MsgType) {
	hdr.t = t
}

func (hdr *MsgHdr) Type() MsgType {
	return hdr.t
}

func (hdr *MsgHdr) String() string {
	return hdr.t.String()
}

type MsgOK struct {
	MsgHdr
}

type MsgError struct {
	MsgHdr
	Message string
}

type MsgAddKey struct {
	MsgHdr
	Data []byte
}

type MsgQuestion struct {
	MsgHdr
	Questions []string
	Echos     []bool
}

type MsgAnswer struct {
	MsgHdr
	Answers []string
}

type MsgListKeys struct {
	MsgHdr
}

type MsgKeys struct {
	MsgHdr
	Keys []KeyInfo
}

type KeyInfo struct {
	Name string
	Type identity.KeyType
	Size int
	ID   string
}

func RPC(conn net.Conn, msg Msg) (Msg, error) {
	err := SendMessage(conn, msg)
	if err != nil {
		return nil, err
	}
	return ReceiveMessage(conn)
}

func SendMessage(conn net.Conn, msg Msg) error {
	var buf [8]byte
	out := new(bytes.Buffer)

	buf[0] = byte(msg.Type())

	_, err := out.Write(buf[:1])
	if err != nil {
		return err
	}

	data, err := encoding.Marshal(msg)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	_, err = out.Write(buf[:4])
	if err != nil {
		return err
	}
	_, err = out.Write(data)
	if err != nil {
		return err
	}

	_, err = conn.Write(out.Bytes())
	return err
}

func ReceiveMessage(conn net.Conn) (Msg, error) {
	var hdr [5]byte

	_, err := io.ReadFull(conn, hdr[:])
	if err != nil {
		return nil, err
	}
	var msg Msg
	msgType := MsgType(hdr[0])

	switch msgType {
	case OK:
		msg = new(MsgOK)

	case Error:
		msg = new(MsgError)

	case AddKey:
		msg = new(MsgAddKey)

	case Question:
		msg = new(MsgQuestion)

	case Answer:
		msg = new(MsgAnswer)

	case ListKeys:
		msg = new(MsgListKeys)

	case Keys:
		msg = new(MsgKeys)

	default:
		return nil, fmt.Errorf("protocol: unexpected message: %s",
			MsgType(hdr[0]))
	}

	len := binary.BigEndian.Uint32(hdr[1:5])
	msgReader := io.LimitedReader{
		R: conn,
		N: int64(len),
	}

	err = encoding.Unmarshal(&msgReader, msg)
	if err != nil {
		return nil, err
	}
	if msgReader.N != 0 {
		return nil, fmt.Errorf("Invalid message: %d bytes extra\n", msgReader.N)
	}
	msg.SetType(msgType)

	return msg, nil
}
