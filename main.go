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

	--version
		Prints out the version tag.
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
	"path"
	"strconv"
)

const BloggyVersionTag = "v0.1-alpha"

// DefaultBlogRepository is the default example blog repository.
const DefaultBlogRepository = "https://github.com/lnsp/bloggy-blueprint"

var (
	// GlobalConfig is the global app configuration.
	Config BlogConfig
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
	versionFlag     = flag.Bool("version", false, "Prints out the version")
)

func initLogger(traceOutput io.Writer) {
	Trace = log.New(traceOutput, "[Trace] ", log.Ltime)
	Info = log.New(os.Stdout, "[Info] ", log.Ltime)
	Warning = log.New(os.Stdout, "[Warning] ", log.Ltime)
	Error = log.New(os.Stderr, "[Error] ", log.Ltime)
}

func Reload() error {
	ClearCache()
	if err := LoadTemplates(); err != nil {
		return err
	}
	if err := LoadPosts(); err != nil {
		return err
	}
	if err := LoadPages(); err != nil {
		return err
	}
	ClearNav()
	AddLinks()
	return nil
}

// Reset resets the blog, deletes all files and clones the source repo.
func Reset(repository string) error {
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
	Info.Println("cloning repository")
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

	if *versionFlag {
		fmt.Println("bloggy", BloggyVersionTag)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		Error.Println("could not determine current working directory")
		return
	}

	BlogFolder = path.Join(cwd, *folderFlag)
	// check if reset
	if *resetFlag {
		Info.Println("reset blog folder")
		err = Reset(*repositoryFlag)
		if err != nil {
			Error.Println("ResetBlog:", err)
			return
		}
	}

	// load configuration file from disk
	err = LoadConfig()
	if err != nil {
		Error.Println("LoadConfig:", err)
		return
	}

	// Load all posts and templates
	err = Reload()
	if err != nil {
		Error.Println("Reload:", err)
		return
	}

	// Start interactive command line interface
	if *interactiveFlag {
		go startInteractiveMode()
	}
	// Create handler
	router := http.Handler(LoadRoutes())
	if *certFlag != "" {
		// Enable HSTS header if HTTPS is enabled
		router = hstsHandler(router)
		Info.Println("enabled TLS/SSL using certificates", *certFlag, "and", *keyFlag)
		go func() {
			address := ":" + strconv.Itoa(Config.Server.TLSPort)
			Error.Println(http.ListenAndServeTLS(address, *certFlag, *keyFlag, router))
		}()
	}
	address := ":" + strconv.Itoa(Config.Server.Port)
	Error.Println(http.ListenAndServe(address, router))
}
