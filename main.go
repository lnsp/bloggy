/*
Bloggy is a fast and simple blogging environment.

Usage:
	bloggy [flags]

Available flags are:

	--reset
		Reset the blog folder and deletes all files.

	--blog="my-blog"
		Set the source folder for templates and posts.

	--repo="https://github.com/lnsp/example-blog"
		Change the git source repository for new blogs
		
	-i
		Enables the interactive command line interface.
		
	-c="certificate"
		Loads the certificate (.key/.pem) and hosts a HTTPS server.
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// DefaultBlogRepository is the default example blog repository.
const DefaultBlogRepository = "https://github.com/lnsp/bloggy-blueprint"

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
	// Logging flags
	logFlags = log.Ltime | log.Lshortfile
)

func runCLI() {
	for {
		fmt.Print("> ")
		var command, arg string
		fmt.Scanln(&command, &arg)

		switch command {
		case "reload":
			// Reload all templates and posts
			LoadTemplates()
			LoadPosts()
			Info.Println("Reloaded posts and templates")
		case "stop":
			// Stops the server
			Info.Println("Forcing shutdown")
			os.Exit(1)
		case "debug":
			switch arg {
			case "on":
				logFlags = log.Ldate | log.Ltime | log.Lshortfile
				initLogger(os.Stdout)
				Info.Println("Activated debug mode")
			case "off":
				logFlags = log.Ldate | log.Ltime
				initLogger(ioutil.Discard)
				Info.Println("Deactivated debug mode")
			}
		}
	}
}

func initLogger(traceOutput io.Writer) {
	Trace = log.New(traceOutput, "[Trace] ", log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[Info] ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[Warning] ", log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[Error] ", log.Ltime|log.Lshortfile)
}

// ResetBlog resets the blog, deletes all files and clones the source repo.
func ResetBlog(repository string) error {
	// delete existing folder
	remError := os.RemoveAll(BlogFolder)
	if remError != nil && !os.IsNotExist(remError) {
		return remError
	}

	cloneCmd := exec.Command("git", "clone", repository, BlogFolder)
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
	initLogger(ioutil.Discard)

	// parse command line arguments
	resetFlag := flag.Bool("reset", false, "Resets the blog data")
	folderFlag := flag.String("blog", "my-blog", "Sets the data source folder")
	repositoryFlag := flag.String("repo", DefaultBlogRepository, "Change the git source repository for resets")
	interactiveFlag := flag.Bool("i", false, "Runs an interactive CLI")
	certificateFlag := flag.String("c", "", "The base certificate path (appends .key and .pem)")
	flag.Parse()

	BlogFolder = *folderFlag
	// check if reset
	if *resetFlag {
		Info.Println("Resetting blog folder:", BlogFolder)
		resetError := ResetBlog(*repositoryFlag)
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

	if *interactiveFlag {
		go runCLI()
	}
	router := http.Handler(LoadRoutes())
	if *certificateFlag != "" {
		certFile := *certificateFlag + ".pem"
		keyFile := *certificateFlag + ".key"
		router = hstsHandler(router)
		Info.Println("Enabled TLS/SSL using certificates", certFile, "and", keyFile)
		go func() {
			err := http.ListenAndServeTLS(GlobalConfig.HostAddressTLS, certFile, keyFile, router)
			if err != nil {
				panic(err)
			}
		}()
	}
	
	err := http.ListenAndServe(GlobalConfig.HostAddress, router)
	if err != nil {
		panic(err)
	}
}
