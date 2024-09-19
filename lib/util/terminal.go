//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package util

import (
	"bytes"
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadPassphrase reads a passphrase from the terminal.
func ReadPassphrase(prompt string, confirm bool) ([]byte, error) {
	for {
		fmt.Printf("%s: ", prompt)
		first, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Printf("\n")
		if err != nil {
			return nil, err
		}
		if !confirm {
			return first, nil
		}

		fmt.Print("Enter same passphrase again: ")
		second, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Printf("\n")
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(first, second) {
			fmt.Print("Passphrases do not match\n")
			continue
		}
		return first, nil
	}
}
