package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	port     = 8443
	endpoint = "/demo"
)

func main() {

	go run_server()

	time.Sleep(2 * time.Second)

	run_client()
}

func run_server() {
	pair, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc(endpoint, handler)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%v", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
	}
	server.Handler = mux

	log.Printf("Server listen at %s", server.Addr)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed start server: %v", err)
	}
}

func run_client() {
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://localhost:%d%s", port, endpoint))
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Println(string(body))
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := []byte("OK")
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
