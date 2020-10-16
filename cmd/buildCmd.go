package cmd

import (
	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(buildCmd)
}

func build(cmd *cobra.Command, args []string) error {
	/*ctx := context.Background()
	srcChan := internal.Browse(ctx, input, threadCount)
	hashedChan := internal.Hash(ctx, srcChan, threadCount)
	internal.ToConsole(ctx, hashedChan)*/
	return nil
}
