//
// traverse.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package local

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/markkurossi/backup/storage"
)

const SpecialMask = os.ModeSymlink | os.ModeDevice | os.ModeNamedPipe |
	os.ModeSocket | os.ModeCharDevice

func Traverse(root string, writer storage.Writer) error {
	fi, err := os.Lstat(root)
	if err != nil {
		return err
	}
	mode := fi.Mode()
	if (mode & SpecialMask) != 0 {
		return nil
	}

	if (mode & os.ModeDir) != 0 {
		files, err := ioutil.ReadDir(root)
		if err != nil {
			return err
		}
		for _, f := range files {
			err = Traverse(fmt.Sprintf("%s/%s", root, f.Name()), writer)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("%s\n", root)
	}

	return nil
}
