//
// Copyright (c) 2018-2024 Markku Rossi
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

// MsgType defines the message types.
type MsgType uint8

// The protocol message types.
const (
	OK        MsgType = 0
	Error             = 1
	AddKey            = 2
	Question          = 3
	Answer            = 4
	ListKeys          = 5
	Keys              = 6
	Decrypt           = 7
	Decrypted         = 8
)

var msgTypeNames = map[MsgType]string{
	OK:        "ok",
	Error:     "error",
	AddKey:    "add-key",
	Question:  "question",
	Answer:    "answer",
	ListKeys:  "list-keys",
	Keys:      "keys",
	Decrypt:   "decrypt",
	Decrypted: "decrypted",
}

func (t MsgType) String() string {
	name, ok := msgTypeNames[t]
	if ok {
		return name
	}
	return fmt.Sprintf("{MsgType %d}", t)
}

// Msg defines the protocol messages.
type Msg interface {
	SetType(t MsgType)
	Type() MsgType
}

// MsgHdr defines the common message header.
type MsgHdr struct {
	t MsgType `backup:"-"`
}

// SetType sets the message type.
func (hdr *MsgHdr) SetType(t MsgType) {
	hdr.t = t
}

// Type returns the message type.
func (hdr *MsgHdr) Type() MsgType {
	return hdr.t
}

func (hdr *MsgHdr) String() string {
	return hdr.t.String()
}

// MsgOK implements the OK message.
type MsgOK struct {
	MsgHdr
}

// MsgError implements the error message.
type MsgError struct {
	MsgHdr
	Message string
}

// MsgAddKey implements the add key message.
type MsgAddKey struct {
	MsgHdr
	Data []byte
}

// MsgQuestion implements the question message.
type MsgQuestion struct {
	MsgHdr
	Questions []string
	Echos     []bool
}

// MsgAnswer implements the answer message.
type MsgAnswer struct {
	MsgHdr
	Answers []string
}

// MsgListKeys implements the list keys message.
type MsgListKeys struct {
	MsgHdr
}

// MsgKeys implements the keys message.
type MsgKeys struct {
	MsgHdr
	Keys []KeyInfo
}

// KeyInfo defines key information.
type KeyInfo struct {
	Name      string
	Type      identity.KeyType
	Size      int
	ID        string
	PublicKey []byte
}

// MsgDecrypt implements the data decryption message.
type MsgDecrypt struct {
	MsgHdr
	KeyID string
	Data  []byte
}

// MsgDecrypted implements the decrypted data message.
type MsgDecrypted struct {
	MsgHdr
	Data []byte
}

// RPC sends the mssage msg to the connection and returns the response
// message.
func RPC(conn net.Conn, msg Msg) (Msg, error) {
	err := SendMessage(conn, msg)
	if err != nil {
		return nil, err
	}
	return ReceiveMessage(conn)
}

// SendMessage sends the message msg to the connection.
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

// ReceiveMessage receives a message from the connection.
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

	case Decrypt:
		msg = new(MsgDecrypt)

	case Decrypted:
		msg = new(MsgDecrypted)

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
		return nil, fmt.Errorf("invalid message: %d bytes extra", msgReader.N)
	}
	msg.SetType(msgType)

	return msg, nil
}
