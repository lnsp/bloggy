package cmd

import "github.com/spf13/cobra"

var Version = "v0.1"

var rootDebug bool

var rootCmd = &cobra.Command{
	Use:     "bloggy",
	Short:   "Minimal blogging engine",
	Version: Version,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootDebug, "debug", "d", false, "Enable debug mode")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
