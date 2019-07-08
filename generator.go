package main

import (
	"errors"
	"sync/atomic"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	errNilTarget = errors.New("Nil target")
)

func newTargeter(conf *Config) (vegeta.Targeter, error) {
	i := int64(-1)
	targets := []*Target{}

	for _, t := range conf.Targets {
		err := t.parseTemplate()
		if err != nil {
			return nil, err
		}

		// Sampling target
		for i := 0; i < t.Repeat; i++ {
			targets = append(targets, t)
		}
	}

	if !conf.Meta.SkipShuffle {
		shuffleRequest(targets)
	}

	return func(t *vegeta.Target) error {
		if t == nil {
			return errNilTarget
		}

		offset := atomic.AddInt64(&i, 1) % int64(len(targets))
		selected := targets[offset]

		if err := selected.toVegetaTarget(t, nil); err != nil {
			return err
		}

		printDebug(t)
		return nil
	}, nil
}
