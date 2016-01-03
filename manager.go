package main

import (
	"github.com/mooxmirror/blog/config"
	"github.com/mooxmirror/blog/templates"
	"github.com/mooxmirror/blog/posts"
	"os/exec"
	"os"
)

const (
	DEFAULT_BLOG_REPOSITORY = "https://github.com/mooxmirror/example-blog"
)

func Reset(folder string) error {
	// delete existing folder
	remError := os.RemoveAll(folder)
	if remError != nil && !os.IsNotExist(remError) {
		return remError
	}

	cloneCmd := exec.Command("git", "clone", DEFAULT_BLOG_REPOSITORY, folder)
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

	// rename folder
	return nil
}

// loads the configuration file from the specified folder
func Load(folder string) (config.Config, error) {
	templates.Load(folder)
	posts.Load(folder)
	return config.GetConfig(folder + "/config.json")
}
