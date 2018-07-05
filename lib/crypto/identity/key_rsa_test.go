//
// key_rsa_test.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package identity

import (
	"testing"
)

func TestRSA(t *testing.T) {
	key, err := NewRSAKey("Test Key", 1024)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}
	data, err := key.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal RSA key: %v", err)
	}

	passphrase := []byte("Hello, world!")

	encrypted, err := Encrypt(data, EncrAES128GCM, passphrase,
		KDFPBKDF24096SHA256)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	decrypted, err := Decrypt(passphrase, encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	_, err = Unmarshal(decrypted)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSA key: %v", err)
	}
}
