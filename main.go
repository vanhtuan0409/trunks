package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	f     = flag.String("target", "targets.yml", "Targets config file path")
	r     = flag.Int("rate", 5, "Request per second to send")
	d     = flag.Int("duration", 0, "Duration to run the request (in seconds)")
	o     = flag.String("output", "", "Output file (default \"stdout\")")
	debug = flag.String("debug", "", "Write debug log to file (default discard)")
	v     = flag.Bool("version", false, "Print version and exit")

	// Set at linking time
	Commit  string
	Date    string
	Version string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getDebugFile() io.Writer {
	if *debug == "" {
		return ioutil.Discard
	}
	f, err := os.OpenFile(*debug, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Warning: Failed to open debug log file. Ignore writing debug log. ERR: %v\n", err)
		return ioutil.Discard
	}
	return f
}

func getOutputFile() (io.Writer, error) {
	if *o == "" {
		return os.Stdout, nil
	}
	f, err := os.OpenFile(*o, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func main() {
	flag.Parse()

	if *v {
		fmt.Printf("Version: %s - Commit: %s - Runtime: %s %s %s - Date: %s\n",
			Version,
			Commit,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
			Date)
		return
	}

	targets, err := parseConfig(*f)
	if err != nil {
		fmt.Printf("Failed to parse target config file. ERR: %v\n", err)
		return
	}

	debugger := getDebugFile()
	if f, ok := debugger.(*os.File); ok {
		defer f.Close()
	}

	out, err := getOutputFile()
	if err != nil {
		fmt.Printf("Failed to open output file. ERR: %v\n", err)
		return
	}
	if f, ok := out.(*os.File); ok {
		defer f.Close()
	}

	rate := vegeta.Rate{Freq: *r, Per: time.Second}
	duration := time.Duration(*d) * time.Second
	targeter, err := newTargeter(targets, debugger)
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

	encoder := vegeta.NewEncoder(out)
	for res := range attacker.Attack(targeter, rate, duration, "Bazookaaa!") {
		encoder.Encode(res)
	}
}
