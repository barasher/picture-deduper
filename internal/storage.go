package internal

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strconv"
)

type Storage struct {
	hashs map[string]uint64
}

func newStorage() *Storage {
	s := Storage{}
	s.hashs = make(map[string]uint64)
	return &s
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
		s.hashs[cur[0]] = h
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
	for p, h := range s.hashs {
		if err = w.Write([]string{p, strconv.FormatUint(h, 10)}) ; err != nil {
			return err
		}
	}
	w.Flush()

	return w.Error()
}
