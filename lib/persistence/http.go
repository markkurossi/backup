//
// http.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package persistence

import (
	"errors"
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

func (h *HTTP) Exists(namespace, key string) (bool, error) {
	req, err := http.NewRequest("HEAD", h.makeURL(namespace, key), nil)
	if err != nil {
		return false, err
	}
	// XXX
	req.Header.Add("js.fetch:mode", "no-cors")
	resp, err := h.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode/100 == 2, nil
}

func (h *HTTP) Get(namespace, key string, flags Flags) ([]byte, error) {
	req, err := http.NewRequest("GET", h.makeURL(namespace, key), nil)
	if err != nil {
		return nil, err
	}
	if (flags & NoCache) != 0 {
		req.Header.Add("Cache-Control", "no-cache")
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (h *HTTP) GetAll(namespace string) (map[string][]byte, error) {
	return nil, errors.New("GetAll not supported for HTTP")
}

func (h *HTTP) Set(namespace, key string, data []byte) error {
	return errors.New("Set not supported for HTTP")
}

func (h *HTTP) makeURL(namespace, key string) string {
	return fmt.Sprintf("%s/%s/%s", h.root, namespace, key)
}
