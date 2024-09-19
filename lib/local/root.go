//
// Copyright (c) 2018-2024 Markku Rossi
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
	// BackupDir defines the name of the backup directory.
	BackupDir = ".backup"
)

// Root implements a storage root.
type Root struct {
	Root string
	Meta string
}

// Mkdir creates a directory into the storage root.
func (root *Root) Mkdir(dir string) error {
	d := fmt.Sprintf("%s/%s", root.Meta, dir)
	return os.MkdirAll(d, 0755)
}

// Exists tests if the key exists in the namespace.
func (root *Root) Exists(namespace, key string) (bool, error) {
	dir := fmt.Sprintf("%s/%s", root.Meta, namespace)
	path := fmt.Sprintf("%s/%s", dir, key)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Get gets the value of the key in the namespace.
func (root *Root) Get(namespace, key string) ([]byte, error) {
	dir := fmt.Sprintf("%s/%s", root.Meta, namespace)
	path := fmt.Sprintf("%s/%s", dir, key)
	return ioutil.ReadFile(path)
}

// Set sets the key-value pair to the namespace.
func (root *Root) Set(namespace, key string, value []byte) error {
	dir := fmt.Sprintf("%s/%s", root.Meta, namespace)
	path := fmt.Sprintf("%s/%s", dir, key)
	return ioutil.WriteFile(path, value, 0644)
}

// GetKeys gets all keys of the namespace.
func (root *Root) GetKeys(namespace string) ([]string, error) {
	dir := fmt.Sprintf("%s/%s", root.Meta, namespace)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		keys = append(keys, fi.Name())
	}
	return keys, nil
}

// GetAll get all items from the namespace.
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

// InitRoot initializes a backup to the directory path.
func InitRoot(path string) (*Root, error) {
	root := newRoot(path)

	info, err := os.Stat(root.Root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root directory '%s' is not directory", path)
	}
	// Is the root already initialized?
	_, err = os.Stat(root.Meta)
	if err == nil {
		return nil, fmt.Errorf("root directory '%s' already initialized", path)
	}

	err = os.Mkdir(root.Meta, 0755)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// OpenRoot opens a backup from the root path.
func OpenRoot(path string) (*Root, error) {
	root := newRoot(path)

	_, err := os.Stat(root.Meta)
	if err != nil {
		return nil, fmt.Errorf("could not access root directory: %s", err)
	}

	return root, nil
}

func newRoot(path string) *Root {
	return &Root{
		Root: path,
		Meta: fmt.Sprintf("%s/%s", path, BackupDir),
	}
}
