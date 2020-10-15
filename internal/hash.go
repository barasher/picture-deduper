package internal

import (
	"context"
	"github.com/corona10/goimagehash"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var allowedExtensionRegexp = regexp.MustCompile("(?i).*\\.jp(e)*g")
const fileLoggingKey = "file"

func Hash(ctx context.Context, inChan chan Entry, goRoutineCount int) chan Entry {
	out := make(chan Entry, goRoutineCount)
	l := log.With().Str("component", "hash").Logger()
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(goRoutineCount)

		for i := 0; i < goRoutineCount; i++ {
			go func(id int) {
				l2 := l.With().Str("goroutine", strconv.Itoa(id)).Logger()
				defer wg.Done()
				var err error
				for {
					select {
					case <-ctx.Done():
						return
					case cur, ok := <-inChan:
						if !ok { // task chan closed
							return
						}
						if !isJpeg(cur.Path) {
							l2.Debug().Str(fileLoggingKey, cur.Path).Msg("Unsupported file type")
							continue
						}
						if cur.Hash, err = hashFile(cur.Path); err != nil {
							l2.Error().Str(fileLoggingKey, cur.Path).Msgf("Error while hashing file: %v", err)
							continue
						}
						out <- cur
					}
				}
			}(i)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func hashFile(path string) (uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		return 0, err
	}
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return 0, err
	}
	return hash.GetHash(), nil
}

func isJpeg(path string) bool {
	return allowedExtensionRegexp.MatchString(path)
}
