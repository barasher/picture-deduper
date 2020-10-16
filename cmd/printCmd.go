package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"picture-deduper/internal"
)

var (
	printCmd = &cobra.Command{
		Use:   "print",
		Short: "prints picture hashs",
		RunE:  print,
	}
)

func init() {
	printCmd.Flags().StringVarP(&input, "dir", "d", "", "Root folder")
	printCmd.Flags().IntVarP(&threadCount, "threadCount", "t", 2, "Thread count")
	printCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(printCmd)
}

func print(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	srcChan := internal.Browse(ctx, input, threadCount)
	hashedChan := internal.Hash(ctx, srcChan, threadCount)
	internal.ToConsole(ctx, hashedChan)
	return nil
}
