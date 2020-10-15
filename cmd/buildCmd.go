package cmd

import (
	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Builds deduplication database",
		RunE:  build,
	}
)

func init() {
	buildCmd.Flags().StringVarP(&input, "dir", "d", "", "Root folder")
	buildCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(buildCmd)
}

func build(cmd *cobra.Command, args []string) error {
	return nil
}