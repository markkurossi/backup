//
// hello.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"syscall/js"
)

func main() {
	ch := make(chan js.Value)

	receiver := func(event js.Value) {
		ch <- event
	}

	init := js.Global().Get("init")
	fmt.Printf("init: %v\n", init)

	cb := js.NewEventCallback(js.PreventDefault|js.StopPropagation, receiver)
	init.Invoke(cb, nil)
	defer cb.Release()

loop:
	for {
		event := <-ch
		fmt.Printf("Received event: %v\n", event)

		key := event.Get("key").String()
		keyCode := event.Get("keyCode").Int()
		fmt.Printf("key=%s, keyCode=%d\n", key, keyCode)
		switch key {
		case "Enter":
			break loop

		case "t":
			test()

		case "p":
			panic("Panic!")
		}
	}
}

func test() {
	fmt.Printf("Generating RSA keypair...\n")
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		fmt.Printf("Failed to generate RSA keypair: %s\n", err)
		return
	}
	fmt.Printf("key: %v\n", key)
	alert := js.Global().Get("alert")
	alert.Invoke("Hello, Wasm!")
}
