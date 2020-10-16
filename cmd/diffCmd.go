package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"picture-deduper/internal"
)

var (
	diffCmd = &cobra.Command{
		Use:   "diff",
		Short: "compares pictures",
		RunE:  diff,
	}
)

func init() {
	diffCmd.Flags().IntVarP(&threadCount, "threadCount", "t", 2, "Thread count")
	rootCmd.AddCommand(diffCmd)
}

func diff(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	input := make(chan internal.Entry, threadCount)
	go func() {
		defer close(input)
		for _, v := range args {
			input <- internal.Entry{Path: v}
		}
	}()
	hashedChan := internal.Hash(ctx, input, threadCount)
	internal.Diff(ctx, hashedChan)
	return nil
}
