package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestBrowse(t *testing.T) {
	var tcs = []struct {
		bufSize int
	}{
		{1},
		{2},
		{3},
		{4},
	}
	for _, tc := range tcs {
		t.Run(strconv.Itoa(tc.bufSize), func(t *testing.T) {
			c := Browse(context.TODO(), "../testdata/sampleDir", 2)
			paths := []string{}
			for cur := range c {
				assert.Nil(t, cur.Err)
				paths = append(paths, cur.Path)
			}
			exp := []string{"../testdata/sampleDir/1.jpg", "../testdata/sampleDir/2.jpg", "../testdata/sampleDir/subDir/3.jpg", "../testdata/sampleDir/subDir/a.txt"}
			assert.Equal(t, len(exp), len(paths))
			assert.Subset(t, paths, exp)
		})
	}
}

func TestBrowse_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := Browse(ctx, "../testdata/sampleDir", 2)
	length := 0
	for _ = range c {
		length++
	}
	assert.Equal(t, 0, length)
}
