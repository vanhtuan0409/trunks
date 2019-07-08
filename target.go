package main

import (
	"bytes"
	"errors"
	"net/http"
	"text/template"

	"github.com/Masterminds/sprig"
	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	errUndefinedTemplate = errors.New("Undefined template")
)

type Target struct {
	URL          string            `yaml:"url"`
	Method       string            `yaml:"method"`
	Boby         string            `yaml:"body"`
	Repeat       int               `yaml:"repeat"`
	Headers      map[string]string `yaml:"headers"`
	header       http.Header
	pathTemplate *template.Template
	bodyTemplate *template.Template
}

func (t *Target) parseTemplate() error {
	pTpl, err := template.New("path").Funcs(sprig.TxtFuncMap()).Parse(t.URL)
	if err != nil {
		return err
	}
	t.pathTemplate = pTpl

	if t.Boby != "" && t.Method != http.MethodGet {
		bTpl, err := template.New("body").Funcs(sprig.TxtFuncMap()).Parse(t.Boby)
		if err != nil {
			return err
		}
		t.bodyTemplate = bTpl
	}

	return nil
}

func (t *Target) interpolatePath(data interface{}) (string, error) {
	if t.pathTemplate == nil {
		return "", errUndefinedTemplate
	}

	var buf bytes.Buffer
	if err := t.pathTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t *Target) interpolateBody(data interface{}) (string, error) {
	if t.bodyTemplate == nil {
		return "", errUndefinedTemplate
	}

	var buf bytes.Buffer
	if err := t.bodyTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t *Target) toVegetaTarget(target *vegeta.Target, data interface{}) error {
	target.Method = t.Method
	target.Header = t.header

	path, err := t.interpolatePath(data)
	if err != nil {
		return err
	}
	target.URL = path

	if t.Boby != "" && t.Method != http.MethodGet {
		body, err := t.interpolateBody(data)
		if err != nil {
			return err
		}
		target.Body = []byte(body)
	}

	return nil
}
