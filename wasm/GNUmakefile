GO1.11:=$(HOME)/work/go/bin/go1.11beta2

all: wasm/hello.wasm httpd

wasm/hello.wasm: wasm/hello.go
	cd wasm; GOOS=js GOARCH=wasm $(GO1.11) build -o $(notdir $@)

httpd: httpd.go
	$(GO1.11) build -o $@
