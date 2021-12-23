package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("NETWORKER_PLAYGROUND_PORT")
	if port == "" {
		log.Fatal("NETWORKER_PLAYGROUND_PORT environment variable is unset")
	}

	r := http.NewServeMux()
	r.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("healthy")
	}))
	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request from %s", r.RemoteAddr)
		fmt.Fprintf(w, "hello from %s", r.Host)
	}))
	p := ":" + port
	log.Printf("https server running on port: %s", p)
	log.Fatal(http.ListenAndServeTLS(p, "./tls/cert.pem", "./tls/key.pem", r))
}
