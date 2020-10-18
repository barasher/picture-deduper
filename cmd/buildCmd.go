package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"picture-deduper/internal"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "build hash database",
		RunE:  build,
	}
)

func init() {
	buildCmd.Flags().IntVarP(&threadCount, "threadCount", "t", 2, "Thread count")
	buildCmd.Flags().StringVarP(&input, "dir", "d", "", "Root folder")
	buildCmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	buildCmd.MarkFlagRequired("output")
	buildCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(buildCmd)
}

func build(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	srcChan := internal.Browse(ctx, input, threadCount)
	hashedChan := internal.Hash(ctx, srcChan, threadCount)
	return internal.AppendToFile(ctx, output, hashedChan)
}
