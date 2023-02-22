# WAGO

The [WebAssembly Gateway Interface (WAGI)][wagi] for Go: a Go+Wazero+WASI Hello World.

This is a simple example of how to implement a tiny HTTP service
using WebAssembly and WASI as the interface between the host server
and the guest WebAssembly module.

The WebAssembly module writes to stdout and reads from stdin, interacts
with environment variables and read command-line arguments; but really,
these interaction with the environment are captured and redirected to
HTTP Requests/Responses.

This project is not intended to be complete, and it is meant only for
demo and educational purposes.

## Usage

Start the server with a Wasm (WASI) module.

    go run main/main.go examples/hello.wasm

Your Wasm module is now responding at http://localhost:8080/

The default example under `examples/hello.wasm` has been compiled from 
`examples/hello.c` using the [wasi-sdk][wasi-sdk]

The default example can be invoked without arguments:

    $ curl 'http://localhost:8080/'
    Hello world

Or you can optionally pass a query string to customize the 'hello' message:

    $ curl 'http://localhost:8080/?name=Wazero!'
    Hello Wazero!

[wagi]: https://github.com/deislabs/wagi/
[wasi-sdk]: https://github.com/WebAssembly/wasi-sdk/
