//
// root.go
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
)

const (
	BackupDir = ".backup"
)

type Root struct {
	Root string
	Meta string
}

func newRoot(path string) *Root {
	return &Root{
		Root: path,
		Meta: fmt.Sprintf("%s/%s", path, BackupDir),
	}
}

func (root *Root) Mkdir(dir string) error {
	d := fmt.Sprintf("%s/%s", root.Meta, dir)
	return os.MkdirAll(d, 0755)
}

func (root *Root) Add(namespace, key string, value []byte) error {
	dir := fmt.Sprintf("%s/%s", root.Meta, namespace)
	path := fmt.Sprintf("%s/%s", dir, key)
	return ioutil.WriteFile(path, value, 0644)
}

func (root *Root) GetAll(namespace string) (map[string][]byte, error) {
	dir := fmt.Sprintf("%s/%s", root.Meta, namespace)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	kv := make(map[string][]byte)

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s/%s", dir, fi.Name())
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		kv[fi.Name()] = data
	}

	return kv, nil
}

func InitRoot(path string) (*Root, error) {
	root := newRoot(path)

	info, err := os.Stat(root.Root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("Root directory '%s' is not directory", path)
	}
	// Is the root already initialized?
	_, err = os.Stat(root.Meta)
	if err == nil {
		return nil, fmt.Errorf("Root directory '%s' already initialized", path)
	}

	err = os.Mkdir(root.Meta, 0755)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func OpenRoot(path string) (*Root, error) {
	root := newRoot(path)

	_, err := os.Stat(root.Meta)
	if err != nil {
		return nil, fmt.Errorf("Could not access root directory: %s", err)
	}

	return root, nil
}
