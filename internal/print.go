package internal

import (
	"context"
	"fmt"
)

func ToConsole(ctx context.Context, inChan chan Entry) {
	for {
		select {
		case <-ctx.Done():
			return
		case cur, ok := <-inChan:
			if !ok { // task chan closed
				return
			}
			fmt.Printf("%v: %v\n", cur.Path, cur.Hash)
		}
	}
}
