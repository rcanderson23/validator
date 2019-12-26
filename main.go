package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validateFuncHandler)
	log.Fatal(http.ListenAndServeTLS(":8443", "/pki/server.crt", "/pki/server.key", mux))
}
