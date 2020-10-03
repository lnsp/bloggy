package cmd

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var initRepository string
var initBlog string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new blog folder",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInit(); err != nil {
			logrus.WithError(err).Fatal("failed to init")
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&initBlog, "blog", "b", "content", "Blog folder to use")
	initCmd.Flags().StringVarP(&initRepository, "repository", "r", "https://github.com/lnsp/bloggy-blueprint", "Repository to clone from")
	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	// delete existing folder
	remError := os.RemoveAll(initBlog)
	if remError != nil && !os.IsNotExist(remError) {
		return remError
	}

	cloneCmd := exec.Command("git", "clone", initRepository, initBlog)
	// run git
	cloneStartError := cloneCmd.Start()
	// failed to run git
	if cloneStartError != nil {
		return cloneStartError
	}
	logrus.WithField("url", initRepository).Info("cloning repository")
	// wait for git to finish
	cloneWaitError := cloneCmd.Wait()
	// git runtime error
	if cloneWaitError != nil {
		return cloneWaitError
	}

	return nil
}
