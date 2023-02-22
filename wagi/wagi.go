package wagi

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/evacchi/wago/wazerox"
	"github.com/tetratelabs/wazero"
	"io"
	"net/http"
	"strings"
)

// A WagiSession is an HTTP session exposed as an io.ReadWriter
// the Evaluate method starts the execution
type WagiSession interface {
	io.ReadWriter
	Evaluate() error
}

// wagiSession implements WagiSession
type wagiSession struct {
	module  *wazerox.RunnableModule
	writer  *wagiWriter
	request *http.Request
	env     map[string]string
	query   []string
}

// wagiWriter implements io.Writer
type wagiWriter struct {
	w               http.ResponseWriter
	headersNewlines int
}

// Write implements io.Writer
//
// It is exposed as stdout to the host.
// The host assumes to be able to write the headers followed by two newlines
// and then the body of the HTTP response.
//
// This implementation really delegates to wagiSession to encapsulate
// some internal state.
func (w *wagiSession) Write(b []byte) (n int, err error) {
	n, err = w.writer.Write(b)
	return
}

// Read implements io.Reader
func (w *wagiSession) Read(p []byte) (n int, err error) {
	n, err = w.request.Body.Read(p)
	return
}

func NewWagiSession(response http.ResponseWriter, request *http.Request, module *wazerox.RunnableModule) WagiSession {
	var env = make(map[string]string)
	var query []string

	// for each header `my-header: value`
	// create an "env var" `HTTP_MY_HEADER=value`
	for k, v := range request.Header {
		upper := strings.ToUpper(k)
		env_var := strings.ReplaceAll(upper, "-", "_")
		prefixed_var := fmt.Sprintf("HTTP_%s", env_var)
		env[prefixed_var] = strings.Join(v, ":")
	}

	// we only need the unparsed pairs a=b, c=d
	query = strings.Split(request.URL.RawQuery, "&")

	return &wagiSession{
		module:  module,
		writer:  &wagiWriter{w: response},
		request: request,
		env:     env,
		query:   query,
	}

}

// Write implements io.Writer
func (w *wagiWriter) Write(b []byte) (int, error) {
	if w.headersNewlines > 2 {
		write, err := w.w.Write(b)
		return write, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	var n int
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			// if we find an empty line, meaning no headers will follow
			w.headersNewlines++
			n += 1 // new line
		} else if w.headersNewlines >= 2 {
			// if we have already reached the end of the headers
			n1, err := w.w.Write([]byte(t))
			if err != nil {
				return n, err
			}
			w.w.Write([]byte{'\n'})
			n += n1 + 1
		} else {
			// reset the new line counter: must be >= 2 consecutive
			w.headersNewlines = 0
			// otherwise "parse" the headers
			// split over : and cleanup any extra space
			ss := strings.Split(t, ":")
			k := strings.TrimSpace(ss[0])
			v := strings.TrimSpace(ss[1])
			w.w.Header().Set(k, v)
			n += len(t)
		}
	}
	return n, nil
}

// Evaluate implements the WagiSession interface
//
// It does the plumbing to setup the Wazero runtime
// so that the guest module may write to stdout,
// read from stdin, use the env variables and args
func (w *wagiSession) Evaluate() (err error) {

	// args is URL Path followed by the pairs of the query string a=b, c=d...
	path := w.request.URL.Path
	args := append([]string{path}, w.query...)

	wconf := wazero.NewModuleConfig().
		WithStdout(w).
		WithStdin(w).
		WithStartFunctions("_start").
		WithArgs(args...)

	for k, v := range w.env {
		wconf.WithEnv(k, v)
	}

	ctx := context.Background()
	m, err := w.module.Instantiate(ctx, wconf)
	m.Close(ctx)
	if err != nil {
		return err
	}
	return nil
}
