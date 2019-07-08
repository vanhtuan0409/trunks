package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	o     = flag.String("output", "", "Output file (default: \"stdout\")")
	debug = flag.String("debug", "", "Write debug log to file (default: discard)")

	debugLogFile io.Writer
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

	if *debug == "" {
		debugLogFile = ioutil.Discard
	} else {
		f, err := os.OpenFile(*debug, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			debugLogFile = ioutil.Discard
			fmt.Printf("Warning: Failed to open debug log file. Ignore writing debug log. ERR: %v\n", err)
		} else {
			defer f.Close()
			debugLogFile = f
		}
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

	var encoder vegeta.Encoder
	if *o == "" {
		encoder = vegeta.NewEncoder(os.Stdout)
	} else {
		f, err := os.OpenFile(*o, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Printf("Failed to open output file. ERR: %v\n", err)
			return
		}
		defer f.Close()
		encoder = vegeta.NewEncoder(f)
	}

	for res := range attacker.Attack(targeter, rate, duration, "Bazookaaa!") {
		encoder.Encode(res)
	}
}
