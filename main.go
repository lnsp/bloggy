package main

import (
	"fmt"
	"os"
	"net/http"
	"./config"
	"github.com/gorilla/mux"
)

var cfg config.Configuration

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "happy new year from %s!", cfg.HostCountry)
}

func main() {
	var err error
	cfg, err = config.GetDefaultConfig()

	if err != nil {
		fmt.Println("Error while loading configuration file")
		os.Exit(1)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", HelloHandler)
	http.ListenAndServe(":8080", router)
}
