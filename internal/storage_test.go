package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad_Nominal(t *testing.T) {
	s, err := Load("../testdata/storage.csv")
	assert.Nil(t, err)
	exp := []Entry{
		{Path: "f1", Hash: 1},
		{Path: "f2", Hash: 2},
		{Path: "f3", Hash: 3},
	}
	assert.ElementsMatch(t, exp, s.hashs)
}

func TestLoad_ParseError(t *testing.T) {
	s, err := Load("../testdata/storageWithParseError.csv")
	assert.Nil(t, err)
	exp := []Entry{
		{Path: "f1", Hash: 1},
		{Path: "f3", Hash: 3},
	}
	assert.Equal(t, exp, s.hashs)

}

func TestLoad_NonExistingFile(t *testing.T) {
	_, err := Load("nonExistingFile.txt")
	assert.NotNil(t, err)
}

func TestSave_Nominal(t *testing.T) {
	s := newStorage()
	s.Add(Entry{Path: "f1", Hash: 1}, Entry{Path: "f2", Hash: 2})

	f, err := ioutil.TempFile("/tmp", "picture-deduper_testSave_Nominal")
	assert.Nil(t, err)
	t.Logf("tempFile: %v", f.Name())
	defer os.Remove(f.Name())

	err = s.Save(f.Name())
	assert.Nil(t, err)

	s2, err := Load(f.Name())
	assert.Nil(t, err)

	assert.ElementsMatch(t, s.hashs, s2.hashs)
}

func checkContainsSimilarity(t *testing.T, sims []Similarity, p1, p2 string, distance int) {
	for _, cur := range sims {
		if ((cur.First.Path == p1 && cur.Second.Path == p2) || (cur.First.Path == p2 && cur.Second.Path == p1)) && cur.Distance == distance {
			return
		}
	}
	assert.Failf(t, "Count not find similarity", "looking for %v/%v/%v in %v", p1, p2, distance, sims)
}

func TestFindSimilarities_0(t *testing.T) {
	s := newStorage()
	s.Add(
		Entry{Path: "f1", Hash: 1},
		Entry{Path: "f2", Hash: 2},
		Entry{Path: "f3", Hash: 1},
		Entry{Path: "f4", Hash: 1})

	sim := s.FindSimilarities(context.TODO(), 0, false, 5)
	sims := []Similarity{}
	for cur := range sim {
		sims = append(sims, cur)
	}

	assert.Equal(t, 3, len(sims))
	checkContainsSimilarity(t, sims, "f1", "f3", 0)
	checkContainsSimilarity(t, sims, "f1", "f4", 0)
	checkContainsSimilarity(t, sims, "f3", "f4", 0)
}

func TestFindSimilarities_1(t *testing.T) {
	s := newStorage()
	s.Add(
		Entry{Path: "f1", Hash: 3}, // b11
		Entry{Path: "f2", Hash: 1}, // b01
		Entry{Path: "f3", Hash: 0}, // b00
		Entry{Path: "f4", Hash: 3}) // b11

	sim := s.FindSimilarities(context.TODO(), 1, false, 5)
	sims := []Similarity{}
	for cur := range sim {
		sims = append(sims, cur)
	}

	assert.Equal(t, 3, len(sims))
	checkContainsSimilarity(t, sims, "f1", "f2", 1)
	checkContainsSimilarity(t, sims, "f2", "f3", 1)
	checkContainsSimilarity(t, sims, "f2", "f4", 1)
}

func TestFindSimilarities_1orLess(t *testing.T) {
	s := newStorage()
	s.Add(
		Entry{Path: "f1", Hash: 3}, // b11
		Entry{Path: "f2", Hash: 1}, // b01
		Entry{Path: "f3", Hash: 0}, // b00
		Entry{Path: "f4", Hash: 3}) // b11

	sim := s.FindSimilarities(context.TODO(), 1, true, 5)
	sims := []Similarity{}
	for cur := range sim {
		sims = append(sims, cur)
	}

	assert.Equal(t, 4, len(sims))
	checkContainsSimilarity(t, sims, "f1", "f2", 1)
	checkContainsSimilarity(t, sims, "f2", "f3", 1)
	checkContainsSimilarity(t, sims, "f2", "f4", 1)
	checkContainsSimilarity(t, sims, "f1", "f4", 0)
}
