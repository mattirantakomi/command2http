package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			println("recovered in flushWriter", r)
		}
	}()

	n, err = fw.w.Write(p)

	if fw.f != nil {
		fw.f.Flush()
	}

	return
}

func handler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			println("recovered in handler", r)
		}
	}()

	fw := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}

	args := strings.Fields(os.Getenv("COMMAND"))
	writers := io.MultiWriter(&fw, os.Stdout)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = writers
	cmd.Stderr = writers

	cmd.Start()
	done := make(chan bool, 1)
	go func() {
		cmd.Wait()
		done <- true
	}()

	select {
	case <-w.(http.CloseNotifier).CloseNotify():
		println("connection closed")
		cmd.Process.Kill()
		return
	case <-done:
		println("command exited")
		fw.f.Flush()
	}

	if cmd.ProcessState.ExitCode() != 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func maxClients(h http.Handler, n int) http.Handler {
	sema := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()

		h.ServeHTTP(w, r)
	})
}

func main() {
	addr := ":8080"
	if port, ok := os.LookupEnv("PORT"); ok {
		addr = ":" + port
	}

	if os.Getenv("COMMAND") == "" {
		panic("COMMAND not set")
	}

	http.Handle("/", maxClients(http.HandlerFunc(handler), 1))

	println("listen " + addr)
	http.ListenAndServe(addr, nil)
}
