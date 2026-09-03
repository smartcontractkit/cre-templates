package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/example/cre-confidential-authorization-template/gateway/internal/httpapi"
	"github.com/example/cre-confidential-authorization-template/gateway/internal/store"
)

func main() {
	addr := env("X402_GATEWAY_ADDR", "127.0.0.1:8402")
	token := os.Getenv("X402_GATEWAY_TOKEN")
	if token == "" {
		log.Fatal("X402_GATEWAY_TOKEN is required")
	}
	path := os.Getenv("X402_STORE_PATH")
	if path == "" {
		log.Fatal("X402_STORE_PATH is required")
	}
	s, err := store.NewPersistent(path)
	if err != nil {
		log.Fatal(err)
	}
	api := httpapi.NewWithTrustedSigner(s, token, os.Getenv("X402_AUTH_TRUSTED_SIGNER"))
	srv := &http.Server{Addr: addr, Handler: api.Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("gateway listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
