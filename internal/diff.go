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
				diff := cur.Hash ^ v.Hash
				fmt.Printf("%v ^ %v: %v\n", cur.Path, v.Path, bits.OnesCount64(diff))
			}
			coll = append(coll, cur)
		}
	}
}
