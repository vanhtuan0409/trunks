package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	f     = flag.String("target", "targets.yml", "Targets config file path")
	r     = flag.Int("rate", 5, "Request per second to send")
	d     = flag.Int("duration", 0, "Duration to run the request (in seconds)")
	debug = flag.Bool("debug", false, "Print debug log")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

	targets, err := parseConfig(*f)
	if err != nil {
		fmt.Printf("Failed to parse target config file. ERR: %v\n", err)
		return
	}

	rate := vegeta.Rate{Freq: *r, Per: time.Second}
	duration := time.Duration(*d) * time.Second
	targeter, err := newTargeter(targets)
	if err != nil {
		fmt.Printf("Failed to initialize targeter. ERR: %v", err)
		return
	}
	attacker := vegeta.NewAttacker()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		attacker.Stop()
	}()

	fmt.Printf("Running load test with rate %d and duration %d\n", *r, *d)

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Bazookaaa!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Average: %s\n", metrics.Latencies.Mean)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Total: %d\n", metrics.Requests)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
	for code, count := range metrics.StatusCodes {
		fmt.Printf("Status code %s: %d\n", code, count)
	}
}
