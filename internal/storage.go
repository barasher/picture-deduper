package internal

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strconv"
)

type Storage struct {
	hashs []Entry
}

func newStorage() *Storage {
	s := Storage{}
	s.hashs = []Entry{}
	return &s
}

func (s *Storage) Add(e ...Entry) {
	s.hashs = append(s.hashs, e...)
}

func Load(file string) (*Storage, error) {
	s := newStorage()
	f, err := os.Open(file)
	if err != nil {
		return s, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		cur, err := r.Read()
		if err == io.EOF {
			break
		}
		h, err := strconv.ParseUint(cur[1], 10, 64)
		if err != nil {
			log.Warn().Msgf("Error while parsing hash '%v' concerning file %v: %v", cur[1], cur[0], err)
			continue
		}
		s.hashs = append(s.hashs, Entry{Path: cur[0], Hash: h})
	}
	return s, nil
}

func (s *Storage) Save(file string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, cur := range s.hashs {
		if err = w.Write([]string{cur.Path, strconv.FormatUint(cur.Hash, 10)}); err != nil {
			return err
		}
	}
	w.Flush()

	return w.Error()
}
