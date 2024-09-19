//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package identity

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/markkurossi/backup/lib/encoding"
	"golang.org/x/crypto/pbkdf2"
)

const (
	magic   = 0x42554944
	version = byte(0)
)

// EncrAlg defines an encryption algorithm.
type EncrAlg int

// Encryption algorithms.
const (
	EncrAES128GCM EncrAlg = 0
)

func (e EncrAlg) String() string {
	switch e {
	case EncrAES128GCM:
		return "AES128-GCM"

	default:
		return fmt.Sprintf("{EncAlg %d}", e)
	}
}

// KeyLen returns the key length of the encryption algorithm e.
func (e EncrAlg) KeyLen() int {
	switch e {
	case EncrAES128GCM:
		return 16

	default:
		panic(fmt.Sprintf("Unknown encryption algorithm %s", e))
	}
}

// KDFAlg defines a key derivation function.
type KDFAlg int

func (k KDFAlg) String() string {
	switch k {
	case KDFPBKDF24096SHA256:
		return "PBKDF2-4096-SHA256"

	default:
		return fmt.Sprintf("{KDFAlg %d}", k)
	}
}

const (
	// KDFPBKDF24096SHA256 defines key derivation function PBKDF with
	// 4096 rounds of SHA-256.
	KDFPBKDF24096SHA256 KDFAlg = 0
)

// EncryptedKey implements an encrypted data blob.
type EncryptedKey struct {
	Magic     uint32
	Version   byte
	Name      string
	Salt      []byte
	KDFAlg    KDFAlg
	EncrAlg   EncrAlg
	Encrypted []byte
}

// Encrypt encrypts the data with the encrAlg and passphrase.
func Encrypt(data []byte, encrAlg EncrAlg, name string,
	passphrase []byte, kdfAlg KDFAlg) ([]byte, error) {

	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	key, err := kdf(passphrase, salt, kdfAlg, encrAlg.KeyLen())
	if err != nil {
		return nil, err
	}
	encrypted, err := encrypt(data, encrAlg, key)
	if err != nil {
		return nil, err
	}

	enc := &EncryptedKey{
		Magic:     magic,
		Version:   version,
		Name:      name,
		Salt:      salt,
		KDFAlg:    kdfAlg,
		EncrAlg:   encrAlg,
		Encrypted: encrypted,
	}

	return encoding.Marshal(enc)
}

// Decrypt decrypts the ciphertext with the passphrase.
func Decrypt(ciphertext, passphrase []byte) ([]byte, error) {
	if len(ciphertext) < 5 {
		return nil, errors.New("Truncated ID key blob")
	}
	if binary.BigEndian.Uint32(ciphertext[:4]) != magic {
		return nil, errors.New("Invalid ID key magic")
	}
	if ciphertext[4] != version {
		return nil, fmt.Errorf("invalid ID key version %d", ciphertext[4])
	}

	enc := new(EncryptedKey)
	err := encoding.Unmarshal(bytes.NewReader(ciphertext), enc)
	if err != nil {
		return nil, err
	}

	// Derive encryption key.
	key, err := kdf(passphrase, enc.Salt, enc.KDFAlg, enc.EncrAlg.KeyLen())
	if err != nil {
		return nil, err
	}

	return decrypt(enc.Encrypted, enc.EncrAlg, key)
}

func kdf(passphrase, salt []byte, alg KDFAlg, keyLen int) ([]byte, error) {
	switch alg {
	case KDFPBKDF24096SHA256:
		return pbkdf2.Key(passphrase, salt, 4096, keyLen, sha256.New), nil

	default:
		return nil, fmt.Errorf("unknown KDF algorithm %s", alg)
	}
}

func encrypt(data []byte, alg EncrAlg, key []byte) ([]byte, error) {
	switch alg {
	case EncrAES128GCM:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		nonceSize := aesgcm.NonceSize()
		nonce := make([]byte, nonceSize)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}
		encrypted := aesgcm.Seal(nil, nonce, data, nil)
		return append(nonce, encrypted...), nil

	default:
		return nil, fmt.Errorf("unknown encryption algorithm %s", alg)
	}
}

func decrypt(data []byte, alg EncrAlg, key []byte) ([]byte, error) {
	switch alg {
	case EncrAES128GCM:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		nonceSize := aesgcm.NonceSize()
		if len(data) < nonceSize {
			return nil, errors.New("Truncated cipher block")
		}
		return aesgcm.Open(nil, data[0:nonceSize], data[nonceSize:], nil)

	default:
		return nil, fmt.Errorf("unknown encryption algorithm %s", alg)
	}
}
