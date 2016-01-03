package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"github.com/mooxmirror/go-blog/config"
	"github.com/mooxmirror/go-blog/routes"
)

var cfg config.Config

func main() {
	var cfgError error
	cfg, cfgError = config.GetDefaultConfig()

	if cfgError != nil {
		log.Fatal("Configuration error: ", cfgError)
		os.Exit(1)
	}

	fmt.Println("Server starts listening on", cfg.HostAddress)
	serverError := http.ListenAndServe(cfg.HostAddress, routes.Setup(cfg))

	if serverError != nil {
		log.Fatal("Server error: ", serverError)
	}
}
