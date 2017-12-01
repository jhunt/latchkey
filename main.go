package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	bind := os.Getenv("LATCHKEY_BIND")
	if bind == "" {
		bind = ":8080"
	}

	mem := make(map[string]string)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(os.Stderr, "access: %s %s\n", req.Method, req.URL.Path)
		key := req.URL.Path

		if req.Method == "PUT" {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, "latchkey failed\n")
				return
			}

			w.WriteHeader(200)
			this := string(b)
			switch this {
			case "same", "changed", "not found":
				w.WriteHeader(400)
				fmt.Fprintf(w, "invalid payload\n")
				return
			}

			if prev, ok := mem[key]; ok && prev == this {
				fmt.Fprintf(w, "same\n")
			} else {
				mem[key] = this
				fmt.Fprintf(w, "changed\n")
			}
			return
		}

		if _, ok := mem[key]; !ok {
			w.WriteHeader(404)
			fmt.Fprintf(w, "not found\n")
			return
		}

		if req.Method == "GET" {
			w.WriteHeader(200)
			fmt.Fprintf(w, "%s\n", mem[key])
			return
		}

		w.WriteHeader(405)
		fmt.Fprintf(w, "method %s not allowed\n", req.Method)
	})

	fmt.Fprintf(os.Stderr, "latchkey starting up on %s\n", bind)
	http.ListenAndServe(bind, nil)
}
