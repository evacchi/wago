run:
	go run main/main.go examples/hello.wasm

all: examples build

build:
	go build ./...

examples: hello.wasm

hello.wasm:
	$(WASI_SDK_PATH)/bin/clang -o examples/hello.wasm examples/hello.c

