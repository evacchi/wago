package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/evacchi/wago/wagi"
	"github.com/evacchi/wago/wazerox"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	if len(os.Args) == 1 {
		fail()
	}

	wasmPath := os.Args[1]

	ctx := context.Background()
	rt := wazero.NewRuntime(ctx)
	defer rt.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)
	module, err := wazerox.ReadModuleFromPath(ctx, rt, wasmPath)
	if err != nil {
		os.Stderr.WriteString(
			fmt.Sprintf("An error occurred while trying to read %s \n", wasmPath))
		fail()
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session := wagi.NewWagiSession(w, r, module)
		err := session.Evaluate()
		if err != nil {
			w.WriteHeader(500)
			log.Fatal(err)
		}
	})

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func fail() {
	os.Stderr.WriteString("Usage: wago <path-to-wasm-binary>\n")
	os.Exit(1)
}
