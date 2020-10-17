package internal

import (
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
	s.Add(Entry{Path:"f1", Hash:1}, Entry{Path:"f2", Hash:2})

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
