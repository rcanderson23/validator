package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	configFile := flag.String("config-file", "config.yaml", "path to config file")
	flag.Parse()
	var config Config
	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validateFuncHandler)
	log.Fatal(http.ListenAndServeTLS(":8443", config.TlsCert, config.TlsKey, mux))
}
