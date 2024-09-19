//
// Copyright (c) 2018-2024 Markku Rossi
//
// All rights reserved.
//

package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Filesystem implements filesystem persistence storage accessor.
type Filesystem struct {
	root string
}

// CreateFilesystem creates a filesystem for the argument diretory.
func CreateFilesystem(root string) (*Filesystem, error) {
	_, err := os.Stat(root)
	if err == nil {
		return nil, fmt.Errorf("root directory '%s' already exists", root)
	}
	err = os.Mkdir(root, 0755)
	if err != nil {
		return nil, err
	}

	return &Filesystem{
		root: root,
	}, nil
}

// OpenFilesystem opens the filesystem from the argument directory.
func OpenFilesystem(root string) (*Filesystem, error) {
	_, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("could not access root directory: %s", err)
	}
	return &Filesystem{
		root: root,
	}, nil
}

// Exists implements Reader.Exists.
func (fs *Filesystem) Exists(namespace, key string) (bool, error) {
	path := fmt.Sprintf("%s/%s/%s", fs.root, namespace, key)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Get implements Reader.Get.
func (fs *Filesystem) Get(namespace, key string, flags Flags) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", fs.root, namespace, key)
	return ioutil.ReadFile(path)
}

// GetAll implements Reader.GetAll.
func (fs *Filesystem) GetAll(namespace string) (map[string][]byte, error) {
	dir := fmt.Sprintf("%s/%s", fs.root, namespace)
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

// Set implements Writer.Set.
func (fs *Filesystem) Set(namespace, key string, value []byte) error {
	dir := fmt.Sprintf("%s/%s", fs.root, namespace)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", dir, key)
	return ioutil.WriteFile(path, value, 0644)
}
