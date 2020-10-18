package internal

import (
	"context"
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"io"
	"math/bits"
	"os"
	"strconv"
	"sync"
)

type Storage struct {
	hashs []Entry
}

// TODO add lock

func newStorage() *Storage {
	s := Storage{}
	s.hashs = []Entry{}
	return &s
}

func (s *Storage) Add(e ...Entry) {
	s.hashs = append(s.hashs, e...)
}

func LoadChan(c chan Entry) *Storage {
	s := newStorage()
	for cur := range c {
		s.Add(cur)
	}
	return s
}

func LoadFile(file string) (*Storage, error) {
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

func (s *Storage) Save(ctx context.Context, file string) error {
	entryChan := make(chan Entry, 1)

	go func() {
		defer close(entryChan)
		for _, v := range s.hashs {
			entryChan <- v
		}
	}()

	return AppendToFile(ctx, file, entryChan)
}

func AppendToFile(ctx context.Context, file string, entryChan chan Entry) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	func() {
		for {
			select {
			case <-ctx.Done():
				return
			case cur, ok := <-entryChan:
				if ! ok {
					return
				}
				if err = w.Write([]string{cur.Path, strconv.FormatUint(cur.Hash, 10)}); err != nil {
					log.Error().Msgf("Error while marshalling hash for %v: %v", cur.Path, err)
				}
			}
		}
	}()
	w.Flush()
	return w.Error()
}

type Similarity struct {
	First    Entry
	Second   Entry
	Distance int
}

func (s *Storage) FindSimilarities(ctx context.Context, distance int, orLess bool, threadCount int) chan Similarity {
	keyChan := make(chan int, threadCount)
	go func() {
		for i := 0; i < len(s.hashs); i++ {
			keyChan <- i
		}
		close(keyChan)
	}()

	simChan := make(chan Similarity, threadCount)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(threadCount)
		for i := 0; i < threadCount; i++ {
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case cur, ok := <-keyChan:
						if !ok {
							return
						}

						if distance == 0 {
							for idx := cur + 1; idx < len(s.hashs); idx++ {
								if s.hashs[cur].Hash == s.hashs[idx].Hash {
									simChan <- Similarity{s.hashs[cur], s.hashs[idx], 0}
								}
							}
						} else {
							for idx := cur + 1; idx < len(s.hashs); idx++ {
								d := bits.OnesCount64(s.hashs[cur].Hash ^ s.hashs[idx].Hash)
								if d == distance || (orLess && d < distance) {
									simChan <- Similarity{s.hashs[cur], s.hashs[idx], d}
								}
							}
						}

					}
				}
			}()
		}
		wg.Wait()
		close(simChan)
	}()

	return simChan
}
