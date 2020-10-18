package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"picture-deduper/internal"
)

var (
	simCmd = &cobra.Command{
		Use:   "similarities",
		Short: "Determine similarities between pictures",
		RunE:  findSim,
	}
)

func init() {
	simCmd.Flags().StringVarP(&input, "db", "", "", "Picture hash database")
	simCmd.Flags().IntVarP(&distance, "distance", "d", 0, "Distance")
	simCmd.Flags().IntVarP(&threadCount, "threadCount", "t", 2, "Thread count")
	simCmd.Flags().BoolVarP(&orLess, "orLess", "l", false, "Include 'lower' distances")
	simCmd.MarkFlagRequired("db")
	simCmd.MarkFlagRequired("distance")
	rootCmd.AddCommand(simCmd)
}

func findSim(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	db, err := internal.LoadFile(input)
	if err != nil {
		return err
	}

	simChan := db.FindSimilarities(ctx, distance, orLess, threadCount)
	for {
		select {
		case <-ctx.Done():
			return nil
		case cur, ok := <-simChan:
			if ! ok {
				return nil
			}
			fmt.Printf("%v,%v,%v\n", cur.Distance, cur.First.Path, cur.Second.Path)
		}
	}
	return nil
}
