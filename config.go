package main

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Meta struct {
	Host        string            `yaml:"host"`
	Headers     map[string]string `yaml:"headers"`
	SkipShuffle bool              `yaml:"skip_shuffle"`
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

func mergeHeader(h1 map[string]string, h2 map[string]string) map[string]string {
	res := make(map[string]string, len(h1)+len(h2))
	for k, v := range h1 {
		res[k] = v
	}
	for k, v := range h2 {
		res[k] = v
	}
	return res
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
	for _, t := range targets.Targets {
		header := convertHeaderMap(mergeHeader(targets.Meta.Headers, t.Headers))
		t.header = header
	}

	return &targets, nil
}
