package internal

import (
	"context"
	"fmt"
	"math/bits"
)

func Diff(ctx context.Context, inChan chan Entry) {
	coll := []Entry{}
	for {
		select {
		case <-ctx.Done():
			return
		case cur, ok := <-inChan:
			if !ok { // task chan closed
				return
			}
			for _, v := range coll {
				fmt.Printf("%v ^ %v: %v\n", cur.Path, v.Path, diff(cur, v))
			}
			coll = append(coll, cur)
		}
	}
}

func diff(e1, e2 Entry) int {
	return bits.OnesCount64(e1.Hash ^ e2.Hash)
}
