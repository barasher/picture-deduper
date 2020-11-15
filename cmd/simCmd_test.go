package cmd

import (
	"github.com/stretchr/testify/assert"
	"picture-deduper/internal"
	"testing"
)

func ExampleExportSimilarityWithoutDetails() {
	sim := internal.Similarity{
		First:    internal.Entry{
			Path: "/a/b/c.jpg",
			Hash: 42,
		},
		Second:   internal.Entry{
			Path: "/d/e/f.jpg",
			Hash: 84,
		},
		Distance: 3,
	}
	exportSimilarityWithoutDetails(sim)
	// Output:
	// 3,/a/b/c.jpg,/d/e/f.jpg
}

func TestExportSimilarityWithDetails_Errors(t *testing.T) {
	var tcs = []struct {
		inId string
		inFirstPath string
		inSecondPath string
		expSuccess bool
	}{
		{"1","../testdata/sampleDir/1.jpg", "../testdata/sampleDir/2.jpg", true},
		{"2","/a/b/c.jpg", "../testdata/sampleDir/1.jpg", false},
		{"3", "../testdata/sampleDir/1.jpg", "/a/b/c.jpg", false},
	}
	for _, tc := range tcs {
		t.Run(tc.inId, func(t *testing.T) {
			sim := internal.Similarity{
				First:    internal.Entry{
					Path: tc.inFirstPath,
					Hash: 42,
				},
				Second:   internal.Entry{
					Path: tc.inSecondPath,
					Hash: 84,
				},
				Distance: 3,
			}
			assert.Equal(t, tc.expSuccess, exportSimilarityWithDetails(sim) == nil)
		})
	}
}

func ExampleExportSimilarityWithDetails() {
	sim := internal.Similarity{
		First:    internal.Entry{
			Path: "../testdata/sampleDir/1.jpg",
			Hash: 42,
		},
		Second:   internal.Entry{
			Path: "../testdata/sampleDir/2.jpg",
			Hash: 84,
		},
		Distance: 3,
	}
	exportSimilarityWithDetails(sim)
	// Output:
	// 3,../testdata/sampleDir/1.jpg,../testdata/sampleDir/2.jpg,77668,77668,1.jpg,2.jpg,../testdata/sampleDir,../testdata/sampleDir,42,84
}