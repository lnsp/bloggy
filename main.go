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
	"io/ioutil"
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
	// Trace logger
	Trace *log.Logger
	// Info logger
	Info *log.Logger
	// Warning logger
	Warning *log.Logger
	// Error logger
	Error *log.Logger
)

func initLogger() {
	Trace = log.New(ioutil.Discard, "[Trace] ", log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[Info] ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[Warning] ", log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[Error] ", log.Ltime|log.Lshortfile)
}

// ResetBlog resets the blog, deletes all files and clones the source repo.
func ResetBlog(folder string) error {
	// delete existing folder
	remError := os.RemoveAll(folder)
	if remError != nil && !os.IsNotExist(remError) {
		return remError
	}

	cloneCmd := exec.Command("git", "clone", DefaultBlogRepository, folder)
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
	resetBlog := flag.Bool("reset", false, "Resets the blog data")
	blogFolder := flag.String("blog", "my-blog", "Sets the data source folder")

	flag.Parse()

	// check if reset
	if *resetBlog {
		Info.Println("Resetting blog folder:", *blogFolder)
		resetError := ResetBlog(*blogFolder)
		if resetError != nil {
			Error.Println("Failed to reset blog:", resetError)
			os.Exit(1)
		}
	}

	// load configuration file from disk
	loadError := LoadConfig(*blogFolder)
	if loadError != nil {
		Error.Println("Failed to load configuration:", loadError)
		os.Exit(1)
	}
	LoadTemplates(*blogFolder)
	LoadPosts(*blogFolder)

	router := LoadRoutes()
	serverError := http.ListenAndServe(GlobalConfig.HostAddress, router)

	if serverError != nil {
		Error.Println("Server error:", serverError)
	}
}
