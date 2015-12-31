package main

import (
	"fmt"
	"os"
	"net/http"
	"./config"
)

var cfg config.Configuration

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "happy new year from %s!", cfg.HostCountry)
}

func main() {
	var err error
	cfg, err = config.GetDefaultConfig()

	if err != nil {
		fmt.Println("Error while loading configuration file")
		os.Exit(1)
	}

	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}
