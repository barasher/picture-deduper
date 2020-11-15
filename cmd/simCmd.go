package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
	simCmd.Flags().BoolVarP(&withDetails, "withDetails", "w", false, "Also export details")
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

	 expFct := exportSimilarityWithoutDetails
	 if withDetails {
		 expFct = exportSimilarityWithDetails
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
			expFct(cur)
		}
	}
	return nil
}

func exportSimilarityWithDetails(sim internal.Similarity) error {
	fStat, err := os.Stat(sim.First.Path)
	if err != nil {
		return err
	}
	sStat, err := os.Stat(sim.Second.Path)
	if err != nil {
		return err
	}
	fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
		sim.Distance,
		sim.First.Path, sim.Second.Path,
		fStat.Size(), sStat.Size(),
		filepath.Base(sim.First.Path), filepath.Base(sim.Second.Path),
		filepath.Dir(sim.First.Path), filepath.Dir(sim.Second.Path),
		sim.First.Hash, sim.Second.Hash)
	return nil
}

func exportSimilarityWithoutDetails(sim internal.Similarity) error {
	fmt.Printf("%v,%v,%v\n", sim.Distance, sim.First.Path, sim.Second.Path)
	return nil
}
