package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {

	pair, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/demo", handler)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%v", 8443),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
	}
	server.Handler = mux

	log.Printf("listen at %s", server.Addr)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed start server: %v", err)
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := []byte("OK")
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
