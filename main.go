package main

import (
	"os"
	"log"
	"flag"
	"net/http"
	"github.com/mooxmirror/blog/routes"
)

func main() {
	// parse command line arguments
	resetBlog := flag.Bool("reset", false, "Resets the blog data")
	blogFolder := flag.String("blog", "my-blog", "Sets the data source folder")

	flag.Parse()

	// check if reset
	if *resetBlog {
		log.Println("Resetting blog folder ", *blogFolder)
		resetError := Reset(*blogFolder)
		if resetError != nil {
			log.Fatal(resetError)
			os.Exit(1)
		}
	}

	// load configuration file from disk
	cfg, loadError := Load(*blogFolder)
	if loadError != nil {
		log.Fatal(loadError)
		os.Exit(1)
	}

	// listen and serve
	log.Println("Server starts listening on", cfg.HostAddress)
	router := routes.Setup(cfg)
	serverError := http.ListenAndServe(cfg.HostAddress, router)

	if serverError != nil {
		log.Fatal(serverError)
	}
}
