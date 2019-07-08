package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	vegeta "github.com/tsenart/vegeta/lib"
)

func printDebug(t *vegeta.Target, debugger io.Writer) {
	fmt.Fprintf(debugger, "%s %s\n", t.Method, t.URL)
	if t.Method != http.MethodGet && len(t.Body) > 0 {
		fmt.Fprintf(debugger, "%s\n", string(t.Body))
	}
}

func shuffleRequest(ts []*Target) {
	rand.Shuffle(len(ts), func(i, j int) { ts[i], ts[j] = ts[j], ts[i] })
}
