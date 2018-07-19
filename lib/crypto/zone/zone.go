//
// zone.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package zone

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"strings"

	"github.com/markkurossi/backup/lib/crypto/identity"
	"github.com/markkurossi/backup/lib/encoding"
	"github.com/markkurossi/backup/lib/local"
	"github.com/markkurossi/backup/lib/storage"
	"github.com/markkurossi/backup/lib/tree"
)

const (
	suite        = AES256CBCHMACSHA256
	rootPointer  = "RootPointer"
	rootDistance = 4096
)

var zoneDirs = []string{
	"identities",
	"objects",
}

type Zone struct {
	Name    string
	Local   *local.Root
	Head    *tree.Snapshot
	HeadID  storage.ID
	idHash  hash.Hash
	secret  []byte
	suite   Suite
	cipher  cipher.Block
	hmac    hash.Hash
	Written uint64
	Saved   uint64
}

func (zone *Zone) identities() string {
	return fmt.Sprintf("%s/identities", zone.Name)
}

func (zone *Zone) objectNames(id storage.ID) (string, string) {
	ns := fmt.Sprintf("%s/objects/%x/%x", zone.Name, id.Data[:1], id.Data[1:2])
	key := fmt.Sprintf("%x", id.Data[2:])

	return ns, key
}

func (zone *Zone) AddIdentity(key identity.PublicKey) error {
	encrypted, err := key.Encrypt(zone.secret)
	if err != nil {
		return err
	}
	return zone.Local.Set(zone.identities(), key.ID(), encrypted)
}

func (zone *Zone) ResolveID(idstring string) (id storage.ID, err error) {
	mid := strings.Index(idstring, "...")
	if mid < 0 {
		// No separator, full ID.
		return storage.IDFromString(idstring)
	}
	if mid < 4 {
		return id, fmt.Errorf("Invalid truncated ID '%s'", idstring)
	}
	prefix, err := hex.DecodeString(idstring[:mid])
	if err != nil {
		return
	}
	suffix, err := hex.DecodeString(idstring[mid+3:])
	if err != nil {
		return
	}

	namespace, _ := zone.objectNames(storage.NewID(prefix))
	keys, err := zone.Local.GetKeys(namespace)
	if err != nil {
		return
	}

	nsBytes := prefix[:2]
	prefix = prefix[2:]

	for _, key := range keys {
		kb, err := hex.DecodeString(key)
		if err != nil {
			fmt.Printf("Skipping invalid object key '%s': %s\n", key, err)
			continue
		}
		if bytes.HasPrefix(kb, prefix) && bytes.HasSuffix(kb, suffix) {
			id.Data = append(id.Data, nsBytes...)
			id.Data = append(id.Data, kb...)
			return id, nil
		}
	}
	return id, fmt.Errorf("Invalid ID '%s'", idstring)
}

// Read ipmlements the storage.Reader interface.
func (zone *Zone) Read(id storage.ID) ([]byte, error) {
	namespace, key := zone.objectNames(id)

	data, err := zone.Local.Get(namespace, key)
	if err != nil {
		return nil, err
	}

	return zone.decrypt(data)
}

// Write implements the storage.Writer interface.
func (zone *Zone) Write(data []byte) (id storage.ID, err error) {
	zone.idHash.Reset()
	zone.idHash.Write(data)

	id = storage.NewID(zone.idHash.Sum(nil))

	namespace, key := zone.objectNames(id)
	err = zone.Local.Mkdir(namespace)
	if err != nil {
		return
	}

	var encrypted []byte
	encrypted, err = zone.encrypt(data)
	if err != nil {
		return
	}

	err = zone.Local.Set(namespace, key, encrypted)
	if err != nil {
		return
	}

	return
}

func (zone *Zone) init(secret []byte, suite Suite) error {
	if len(secret) != suite.KeyLen() {
		return fmt.Errorf("Invalid zone key length: %d vs %d", len(secret),
			zone.suite.KeyLen())
	}
	zone.secret = secret
	zone.suite = suite

	split1 := suite.IDHashKeyLen()
	split2 := split1 + suite.CipherKeyLen()

	switch suite {
	case AES256CBCHMACSHA256:
		zone.idHash = hmac.New(sha256.New, secret[:split1])

		block, err := aes.NewCipher(secret[split1:split2])
		if err != nil {
			return err
		}
		zone.cipher = block

		zone.hmac = hmac.New(sha256.New, secret[split2:])

	default:
		return fmt.Errorf("Unsupported suite: %s", suite)
	}

	return nil
}

func (zone *Zone) getHead() error {
	data, err := zone.Local.Get(zone.Name, rootPointer)
	if err != nil {
		return zone.bruteForceRootPointer()
	}
	in := bytes.NewReader(data)

	ptr1 := new(RootPointer)
	err1 := encoding.Unmarshal(in, ptr1)

	ptr2 := new(RootPointer)
	var err2 error

	if len(data) > rootDistance {
		in = bytes.NewReader(data[rootDistance:])
		err2 = encoding.Unmarshal(in, ptr2)
	} else {
		err2 = io.EOF
	}

	if err1 == nil {
		err1 = zone.checkRootPointer(ptr1)
	}
	if err2 == nil {
		err2 = zone.checkRootPointer(ptr2)
	}

	var id storage.ID

	if err1 == nil && err2 == nil {
		if ptr1.Timestamp > ptr2.Timestamp {
			id = ptr1.Pointer
		} else {
			id = ptr2.Pointer
		}
	} else if err1 == nil {
		id = ptr1.Pointer
	} else if err2 == nil {
		id = ptr2.Pointer
	} else {
		return zone.bruteForceRootPointer()
	}

	element, err := tree.DeserializeID(id, zone)
	if err != nil {
		return zone.bruteForceRootPointer()
	}
	head, ok := element.(*tree.Snapshot)
	if !ok {
		return fmt.Errorf("Root is not a snapshot (%T)", element)
	}

	zone.Head = head
	zone.HeadID = id

	return nil
}

func (zone *Zone) checkRootPointer(ptr *RootPointer) error {
	return fmt.Errorf("checkRootPointer not implemented yet")
}

func (zone *Zone) bruteForceRootPointer() error {
	var best *tree.Snapshot
	var bestID storage.ID
	var buf [2]byte

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			buf[0] = byte(i)
			buf[1] = byte(j)

			ns, _ := zone.objectNames(storage.NewID(buf[:]))
			kvs, err := zone.Local.GetAll(ns)
			if err != nil {
				continue
			}

			for k, v := range kvs {
				data, err := zone.decrypt(v)
				if err != nil {
					continue
				}
				element, err := tree.Deserialize(data, zone)
				if err != nil {
					continue
				}
				snapshot, ok := element.(*tree.Snapshot)
				if ok {
					if best == nil || snapshot.Timestamp > best.Timestamp {

						idData := []byte{byte(i), byte(j)}
						suffix, err := hex.DecodeString(k)
						if err != nil {
							continue
						}
						idData = append(idData, suffix...)

						best = snapshot
						bestID = storage.NewID(idData)
					}
				}
			}
		}
	}

	if best == nil {
		return fmt.Errorf("No root pointer found from object store")
	}

	zone.Head = best
	zone.HeadID = bestID

	return nil
}

func (zone *Zone) encrypt(orig []byte) ([]byte, error) {
	zone.Written += uint64(len(orig))

	// Does it compress?
	var b bytes.Buffer
	z := zlib.NewWriter(&b)
	z.Write(orig)
	z.Close()

	compressed := b.Bytes()
	var data []byte
	if len(compressed) < len(orig) {
		zone.Saved += uint64(len(orig) - len(compressed))
		data = append(data, 1)
		data = append(data, compressed...)
	} else {
		data = append(data, 0)
		data = append(data, orig...)
	}

	blockSize := zone.cipher.BlockSize()

	var padLen = blockSize - (len(data) % blockSize)

	inputLen := blockSize + len(data) + padLen + zone.hmac.Size()
	input := make([]byte, blockSize, inputLen)

	// IV
	_, err := io.ReadFull(rand.Reader, input)
	if err != nil {
		return nil, err
	}
	// Data
	input = append(input, data...)

	// Padding.
	for i := 0; i < padLen-1; i++ {
		input = append(input, byte(i))
	}
	// Padding length.
	input = append(input, byte(padLen))

	// Encrypt input
	cbc := cipher.NewCBCEncrypter(zone.cipher, input[:blockSize])
	toCrypt := input[blockSize : blockSize+len(data)+padLen]
	cbc.CryptBlocks(toCrypt, toCrypt)

	// Compute HMAC.
	zone.hmac.Reset()
	zone.hmac.Write(input)

	// Append HMAC to input and return the updated slice.
	return zone.hmac.Sum(input), nil
}

func (zone *Zone) decrypt(data []byte) ([]byte, error) {
	// Sanity check input length.
	blockSize := zone.cipher.BlockSize()
	hmacLen := zone.hmac.Size()
	if len(data) <= blockSize+hmacLen {
		// Zero-length data is impossible because of minimum padding
		// up to next block size (+1 for padding length).
		return nil, fmt.Errorf("Encrypted data too short")
	}
	if (len(data)-hmacLen)%blockSize != 0 {
		// Encrypted data not rounded up to block size.
		return nil, fmt.Errorf("Invalid encrypted data length")
	}
	split := len(data) - hmacLen
	encrypted := data[:split]
	hmac := data[split:]

	// Check HMAC.
	zone.hmac.Reset()
	zone.hmac.Write(encrypted)
	computed := zone.hmac.Sum(nil)
	if !bytes.Equal(hmac, computed) {
		return nil, fmt.Errorf("HMAC mismatch")
	}

	// Decrypt data.
	cbc := cipher.NewCBCDecrypter(zone.cipher, encrypted[:blockSize])
	toDecrypt := encrypted[blockSize:]
	cbc.CryptBlocks(toDecrypt, toDecrypt)

	padLen := int(toDecrypt[len(toDecrypt)-1])
	if padLen > len(toDecrypt) {
		return nil, fmt.Errorf("Invalid padding")
	}

	decrypted := toDecrypt[:len(toDecrypt)-padLen]
	if len(decrypted) == 0 {
		return nil, fmt.Errorf("Truncated data")
	}

	// Was the data compressed?
	if decrypted[0] == 0 {
		// No.
		return decrypted[1:], nil
	}

	r, err := zlib.NewReader(bytes.NewReader(decrypted[1:]))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

func newZone(name string, local *local.Root) *Zone {
	return &Zone{
		Name:  name,
		Local: local,
	}
}

func Create(local *local.Root, name string) (*Zone, error) {
	secret := make([]byte, suite.KeyLen())
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		return nil, err
	}

	for _, dir := range zoneDirs {
		err := local.Mkdir(fmt.Sprintf("%s/%s", name, dir))
		if err != nil {
			return nil, err
		}
	}

	zone := newZone(name, local)
	if err := zone.init(secret, suite); err != nil {
		return nil, err
	}

	return zone, nil
}

func Open(local *local.Root, name string, keys []identity.PrivateKey) (
	*Zone, error) {

	zone := newZone(name, local)

	// Get zone identities.
	identities, err := local.GetAll(zone.identities())
	if err != nil {
		return nil, err
	}

	// Do we have an identity to open the zone?
	for _, key := range keys {
		data, ok := identities[key.ID()]
		if !ok {
			continue
		}
		secret, err := key.Decrypt(data)
		if err != nil {
			continue
		}
		err = zone.init(secret, suite)
		if err != nil {
			return nil, err
		}

		// Get head snapshot.
		err = zone.getHead()
		if err != nil {
			return nil, err
		}

		return zone, nil
	}
	return nil, fmt.Errorf("No key to open zone '%s'", name)
}
