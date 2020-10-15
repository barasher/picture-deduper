package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

const refHash uint64 = 12286558942329960953

func TestHash2(t *testing.T) {
	in := make(chan Entry, 20)
	in <- Entry{Path: "../testdata/sampleDir/1.jpg", Err: nil}
	in <- Entry{Path: "../testdata/sampleDir/2.jpg", Err: nil}
	in <- Entry{Path: "../testdata/sampleDir/subDir/3.jpg", Err: nil}
	close(in)

	out := Hash(context.Background(), in, 1)
	paths := []string{}
	for cur := range out {
		assert.Nil(t, cur.Err)
		paths = append(paths, cur.Path)
	}
	assert.Equal(t, 3, len(paths))
	exp := []string{"../testdata/sampleDir/1.jpg", "../testdata/sampleDir/2.jpg", "../testdata/sampleDir/subDir/3.jpg"}
	assert.Subset(t, paths, exp)
}

func TestIsJpeg(t *testing.T) {
	var tcs = []struct {
		path      string
		expResult bool
	}{
		{"a.jpeg", true},
		{"a.jpg", true},
		{"a.JpG", true},
		{"a.pdf", false},
		{"aa", false},
	}
	for _, tc := range tcs {
		t.Run(tc.path, func(t *testing.T) {
			assert.Equal(t, tc.expResult, isJpeg(tc.path))
		})
	}
}

func TestHashFile(t *testing.T) {
	var tcs = []struct {
		path    string
		expOk   bool
		expHash uint64
	}{
		{"../testdata/sampleDir/1.jpg", true, refHash},
		{"nonExistingFile", false, 0},
		{"../testdata/sampleDir/subDir/a.txt", false, 0},
	}
	for _, tc := range tcs {
		t.Run(tc.path, func(t *testing.T) {
			hash, err := hashFile(tc.path)
			assert.Equal(t, tc.expOk, err == nil)
			if tc.expOk {
				assert.Equal(t, tc.expHash, hash)
			}
		})
	}
}

func TestHash(t *testing.T) {
	in := make(chan Entry, 20)
	in <- Entry{Path: "../testdata/sampleDir/1.jpg"}
	in <- Entry{Path: "../testdata/sampleDir/subDir/a.txt"}
	in <- Entry{Path: "../testdata/sampleDir/2.jpg"}
	in <- Entry{Path: "nonExisting.jpg"}
	in <- Entry{Path: "../testdata/sampleDir/subDir/3.jpg"}
	close(in)

	out := Hash(context.Background(), in, 2)
	res := make(map[string]uint64)
	for cur := range out {
		assert.Nil(t, cur.Err)
		res[cur.Path] = cur.Hash
	}
	assert.Equal(t, 3, len(res))
	assert.Equal(t, refHash, res["../testdata/sampleDir/1.jpg"])
	assert.Equal(t, refHash, res["../testdata/sampleDir/2.jpg"])
	assert.Equal(t, refHash, res["../testdata/sampleDir/subDir/3.jpg"])
}
