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
	fileInfo, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}
	mode := fileInfo.Mode()
	if (mode & SpecialMask) != 0 {
		return nil, nil
	}

	if (mode & os.ModeDir) != 0 {
		files, err := ioutil.ReadDir(root)
		if err != nil {
			return nil, err
		}

		dir := tree.NewDirectory()

		for _, f := range files {
			id, err := Traverse(fmt.Sprintf("%s/%s", root, f.Name()), writer)
			if err != nil {
				return nil, err
			}
			if id == nil {
				// Unsupported file type.
				continue
			}

			if f.IsDir() {
				fmt.Printf("%s\t%s/\n", id, f.Name())
			} else {
				fmt.Printf("%s\t%s\n", id, f.Name())
			}

			dir.Add(f.Name(), uint32(f.Mode()), f.ModTime().Unix(), id)
		}

		data, err := dir.Serialize()
		if err != nil {
			return nil, err
		}
		return writer.Write(data)
	} else {
		if fileInfo.Size() < 1024*1024 {
			data, err := ioutil.ReadFile(root)
			if err != nil {
				return nil, err
			}
			file := tree.NewFile(data)
			data, err = file.Serialize()
			if err != nil {
				return nil, err
			}
			return writer.Write(data)
		} else {
			return nil, fmt.Errorf("File %s > 1MB", fileInfo.Name())
		}
	}
}
