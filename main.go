package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validateFuncHandler)
	//server := &http.Server{
	//	Addr:    ":8443",
	//	Handler: mux,
	//	TLSConfig: &tls.Config{
	//		kk
	//	}
	//}
	log.Fatal(http.ListenAndServeTLS(":8443", "/pki/server.crt", "/pki/server.key", nil))
}
