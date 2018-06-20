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
	"github.com/markkurossi/backup/tree"
)

const SpecialMask = os.ModeSymlink | os.ModeDevice | os.ModeNamedPipe |
	os.ModeSocket | os.ModeCharDevice

func Traverse(root string, writer storage.Writer) (*tree.ID, error) {
	fi, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}
	mode := fi.Mode()
	if (mode & SpecialMask) != 0 {
		return nil, nil
	}

	if (mode & os.ModeDir) != 0 {
		files, err := ioutil.ReadDir(root)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			id, err := Traverse(fmt.Sprintf("%s/%s", root, f.Name()), writer)
			if err != nil {
				return nil, err
			}
			fmt.Printf("%s\t%s\n", id, f.Name())
		}
	} else {
		// XXX Should do in 1MB chunks
		data, err := ioutil.ReadFile(root)
		if err != nil {
			return nil, err
		}
		return writer.Write(data)
	}

	return nil, nil
}
