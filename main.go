/*
Bloggy is a fast and simple blogging environment.

Usage:
	bloggy [flags]

Available flags are:

	--reset
		Reset the blog folder and deletes all files.

	--blog="my-blog"
		Set the source folder for templates and posts.

	--repo="https://github.com/lnsp/bloggy-blueprint"
		Change the git source repository for new blogs

	-i
		Enables the interactive command line interface.
		Type 'help' in the CLI to get more information.

	-c="certificate"
		Loads the certificate file and enables HTTPS.

	-k="key"
		Loads the private key file. Required if HTTPS is enabled.
*/
package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

// DefaultBlogRepository is the default example blog repository.
const DefaultBlogRepository = "https://github.com/lnsp/bloggy-blueprint"

var (
	// GlobalConfig is the global app configuration.
	GlobalConfig BlogConfig
	// BlogFolder is the source folder.
	BlogFolder string
	// Loggers
	Trace, Info, Warning, Error *log.Logger
	// Logging flags
	logFlags = log.Ltime | log.Lshortfile
	// Other flags
	resetFlag       = flag.Bool("reset", false, "Resets the blog data")
	folderFlag      = flag.String("blog", "my-blog", "Sets the data source folder")
	repositoryFlag  = flag.String("repo", DefaultBlogRepository, "Change the git source repository for resets")
	interactiveFlag = flag.Bool("i", false, "Runs an interactive CLI")
	certFlag        = flag.String("c", "", "Certificate file for HTTPS")
	keyFlag         = flag.String("k", "", "Private key file for HTTPS")
)

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
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		Error.Println("Could not determine working directory")
		os.Exit(1)
	}
	BlogFolder = path.Join(cwd, *folderFlag)
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
	// Load all posts and templates
	LoadTemplates()
	LoadPosts()
	LoadPages()

	// Start interactive command line interface
	if *interactiveFlag {
		go runCLI()
	}
	// Create handler
	router := http.Handler(LoadRoutes())
	if *certFlag != "" {
		// Enable HSTS header if HTTPS is enabled
		router = hstsHandler(router)
		Info.Println("Enabled TLS/SSL using certificates", *certFlag, "and", *keyFlag)
		go func() {
			Error.Println(http.ListenAndServeTLS(GlobalConfig.HostAddressTLS, *certFlag, *keyFlag, router))
		}()
	}
	Error.Println(http.ListenAndServe(GlobalConfig.HostAddress, router))
}
