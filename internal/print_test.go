package internal

import (
	"context"
)

func ExampleToConsole() {
	in := make(chan Entry, 20)
	in <- Entry{Path: "turlututu", Hash: 42}
	in <- Entry{Path: "blabla", Hash: 99}
	close(in)
	ToConsole(context.TODO(), in)
	// Output:
	// turlututu: 42
	// blabla: 99
}
