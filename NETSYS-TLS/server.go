package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "Hello, world!")
	case http.MethodPost:
		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			var msg Message
			body, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(body, &msg)
			fmt.Fprintf(w, "Received JSON: %s\n", msg.Text)
		} else {
			r.ParseMultipartForm(10 << 20) // 10 MB
			file, handler, err := r.FormFile("file")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			fmt.Fprintf(w, "Uploaded File: %s\n", handler.Filename)
			fmt.Fprintf(w, "File Size: %d\n", handler.Size)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is listening on port 8080...")

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Error loading certificates:", err)
		return
	}

	server := &http.Server{
		Addr: ":8080",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}