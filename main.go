/*
Bloggy is a fast and simple blogging environment.

Usage:
	bloggy [flags]

Available flags are:

	--reset
		Resets the blog folder and deletes all files.

	--blog="my-blog"
		Sets the source folder for templates and posts.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// DefaultBlogRepository is the default example blog repository.
const DefaultBlogRepository = "https://github.com/mooxmirror/example-blog"

var (
	// GlobalConfig is the global app configuration.
	GlobalConfig BlogConfig
	// BlogFolder is the source folder.
	BlogFolder string
	// Trace logger
	Trace *log.Logger
	// Info logger
	Info *log.Logger
	// Warning logger
	Warning *log.Logger
	// Error logger
	Error *log.Logger
)

func runCLI() {
	for {
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)

		switch input {
		case "reload":
			// Reload all templates and posts
			LoadTemplates()
			LoadPosts()
		case "stop":
			// Stops the server
			Info.Println("Forcing shutdown")
			os.Exit(1)
		}
	}
}

func initLogger() {
	Trace = log.New(os.Stdout, "[Trace] ", log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[Info] ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[Warning] ", log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[Error] ", log.Ltime|log.Lshortfile)
}

// ResetBlog resets the blog, deletes all files and clones the source repo.
func ResetBlog() error {
	// delete existing folder
	remError := os.RemoveAll(BlogFolder)
	if remError != nil && !os.IsNotExist(remError) {
		return remError
	}

	cloneCmd := exec.Command("git", "clone", DefaultBlogRepository, BlogFolder)
	// run git
	cloneStartError := cloneCmd.Start()
	// failed to run git
	if cloneStartError != nil {
		return cloneStartError
	}
	// wait for git to finish
	cloneWaitError := cloneCmd.Wait()
	// git runtime error
	if cloneWaitError != nil {
		return cloneWaitError
	}

	return nil
}

func main() {
	initLogger()

	// parse command line arguments
	resetFlag := flag.Bool("reset", false, "Resets the blog data")
	folderFlag := flag.String("blog", "my-blog", "Sets the data source folder")

	flag.Parse()

	BlogFolder = *folderFlag
	// check if reset
	if *resetFlag {
		Info.Println("Resetting blog folder:", BlogFolder)
		resetError := ResetBlog()
		if resetError != nil {
			Error.Println("Failed to reset blog:", resetError)
			os.Exit(1)
		}
	}

	// load configuration file from disk
	loadError := LoadConfig()
	if loadError != nil {
		Error.Println("Failed to load configuration:", loadError)
		os.Exit(1)
	}
	LoadTemplates()
	LoadPosts()

	go runCLI()
	router := LoadRoutes()
	serverError := http.ListenAndServe(GlobalConfig.HostAddress, router)

	if serverError != nil {
		Error.Println("Server error:", serverError)
	}
}
