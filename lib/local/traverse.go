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
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/markkurossi/backup/lib/storage"
	"github.com/markkurossi/backup/lib/tree"
)

const SpecialMask = os.ModeSymlink | os.ModeDevice | os.ModeNamedPipe |
	os.ModeSocket | os.ModeCharDevice

var ignores = map[string]string{
	".backup":   "Backup info directory",
	".DS_Store": "macOS Desktop Services Store",
}
var ignoreSuffixes = []string{
	"~",
}

func Traverse(root string, writer storage.Writer) (id storage.ID, err error) {
	fileInfo, err := os.Lstat(root)
	if err != nil {
		return
	}
	mode := fileInfo.Mode()
	if (mode & SpecialMask) != 0 {
		return
	}

	// Check system ignores.
	name := fileInfo.Name()
	_, ok := ignores[name]
	if ok {
		return
	}
	for _, suffix := range ignoreSuffixes {
		if strings.HasSuffix(name, suffix) {
			return
		}
	}

	var data []byte

	if (mode & os.ModeDir) != 0 {
		var files []os.FileInfo
		files, err = ioutil.ReadDir(root)
		if err != nil {
			return
		}

		dir := tree.NewDirectory()

		for _, f := range files {
			id, err = Traverse(fmt.Sprintf("%s/%s", root, f.Name()), writer)
			if err != nil {
				return
			}
			if id.Undefined() {
				// Unsupported file type.
				continue
			}

			if f.IsDir() {
				fmt.Printf("%s\t%s/\n", id, f.Name())
			} else {
				fmt.Printf("%s\t%s\n", id, f.Name())
			}

			dir.Add(f.Name(), f.Mode(), f.ModTime().Unix(), id)
		}

		data, err = dir.Serialize()
		if err != nil {
			return
		}
		return writer.Write(data)
	} else {
		if fileInfo.Size() < 1024*1024 {
			data, err = ioutil.ReadFile(root)
			if err != nil {
				return
			}
			file := tree.NewSimpleFile(data)
			data, err = file.Serialize()
			if err != nil {
				return
			}
			return writer.Write(data)
		} else {
			var file *os.File
			file, err = os.Open(root)
			if err != nil {
				return
			}
			defer file.Close()

			buf := make([]byte, 1024*1024)
			cf := tree.NewChunkedFile(fileInfo.Size())

			for {
				var read int

				read, err = file.Read(buf)
				if read == 0 {
					if err != io.EOF {
						return
					}
					break
				}
				id, err = writer.Write(buf[:read])
				if err != nil {
					return
				}
				cf.Add(int64(read), id)
			}

			data, err = cf.Serialize()
			if err != nil {
				return
			}
			return writer.Write(data)
		}
	}
}
