package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	tplFns template.FuncMap = map[string]interface{}{
		"randInt":   randInt,
		"timestamp": timestamp,
	}
)

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func timestamp() int64 {
	return time.Now().Unix()
}

func printDebug(t *vegeta.Target) {
	fmt.Printf("%s %s\n", t.Method, t.URL)
	if t.Method != http.MethodGet && len(t.Body) > 0 {
		fmt.Printf("%s\n", string(t.Body))
	}
}

func shuffleRequest(ts []*Target) {
	rand.Shuffle(len(ts), func(i, j int) { ts[i], ts[j] = ts[j], ts[i] })
}
