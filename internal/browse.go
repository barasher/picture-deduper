package internal

import (
	"context"
	"os"
	"path/filepath"
)

func Browse(ctx context.Context, src string, bufSize int) chan Entry {
	c := make(chan Entry, bufSize)
	go func() {
		defer close(c)
		var e error = nil
		filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				e = filepath.SkipDir
			default:
				if ! info.IsDir() {
					c <- Entry{Path: path, Err: err}
				}
			}
			return e
		})
	}()
	return c
}
