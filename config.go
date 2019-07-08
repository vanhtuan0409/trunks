package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Meta struct {
	Host    string            `yaml:"host"`
	Headers map[string]string `yaml:"headers"`
}

type Config struct {
	Meta    Meta      `yaml:"meta"`
	Targets []*Target `yaml:"targets"`
}

func convertHeaderMap(hm map[string]string) http.Header {
	header := map[string][]string{}
	for key, value := range hm {
		header[key] = []string{value}
	}
	return header
}

func parseConfig(fPath string) (*Config, error) {
	content, err := ioutil.ReadFile(fPath)
	if err != nil {
		return nil, err
	}

	var targets Config
	err = yaml.Unmarshal(content, &targets)
	if err != nil {
		return nil, err
	}

	// Pre-processing
	header := convertHeaderMap(targets.Meta.Headers)
	for _, t := range targets.Targets {
		t.Header = header
		t.Path = fmt.Sprintf("%s%s", targets.Meta.Host, t.Path)
		if t.Repeat < 1 {
			t.Repeat = 1
		}
	}

	return &targets, nil
}
