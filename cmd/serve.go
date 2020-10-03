package cmd

import (
	"fmt"

	"github.com/lnsp/bloggy/pkg/config"
	"github.com/lnsp/bloggy/pkg/content"
	"github.com/lnsp/bloggy/pkg/routes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serveBlog string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a HTTP server and serve the blog content",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runServe(); err != nil {
			logrus.WithError(err).Fatal("failed to serve")
		}
	},
}

func runServe() error {
	// Open config
	cfg, err := config.Load(serveBlog)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	// Create resolver
	resolver := routes.NewResolver(cfg)
	// Open indexer
	index, err := content.NewIndex(cfg, resolver)
	if err != nil {
		return fmt.Errorf("new index: %w", err)
	}
	// Open templater
	templater, err := content.NewTemplater(cfg, index)
	if err != nil {
		return fmt.Errorf("new templater: %w", err)
	}
	// Open server
	router := routes.NewRouter(cfg, templater)
	// Wait and listen
	return router.Serve()
}

func init() {
	serveCmd.Flags().StringVarP(&serveBlog, "blog", "b", "content", "Blog folder to serve")
	rootCmd.AddCommand(serveCmd)
}
