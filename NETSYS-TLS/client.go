package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		fmt.Println("Error loading CA certificate:", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
			DisableKeepAlives: true,
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("GET Response Status:", resp.Status)

	postJSON(client)

	postMultipartForm(client)
}

func postJSON(client *http.Client) {
	msg := map[string]string{"text": "Hello, server!"}
	jsonData, _ := json.Marshal(msg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", "https://localhost:8080", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("POST JSON Response Status:", resp.Status)
}

func postMultipartForm(client *http.Client) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "example.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	writer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", "https://localhost:8080", body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("POST Multipart Form Response Status:", resp.Status)
}