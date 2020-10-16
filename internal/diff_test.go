package internal

import (
	"context"
)

func ExampleDiff() {
	in := make(chan Entry, 20)
	in <- Entry{Path: "turlututu", Hash: 27}
	in <- Entry{Path: "blabla", Hash: 20}
	close(in)
	Diff(context.TODO(), in)
	// Output:
	// blabla ^ turlututu: 4
}
