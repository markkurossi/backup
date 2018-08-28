//
// http.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTP struct {
	root   string
	client *http.Client
}

func NewHTTP(root string) (*HTTP, error) {
	return &HTTP{
		root:   root,
		client: &http.Client{},
	}, nil
}

func (h *HTTP) Read(id ID) ([]byte, error) {
	url, err := h.makeURL(id)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (h *HTTP) makeURL(id ID) (string, error) {
	if len(id.Data) < 2 {
		return "", fmt.Errorf("Invalid ID: %s", id)
	}
	return fmt.Sprintf("%s/objects/%x/%x/%x",
		h.root, id.Data[:1], id.Data[1:2], id.Data[2:]), nil
}
