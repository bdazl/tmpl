//	  _    __  __ ,___  _
//	 | |  /  \/  \  _ \| |
//	<   > | |\/| | |_) | |
//	 | |__| |  | |  __/| |__
//	 \___/|_|  |_|_|   \___/
//
// Copyright (C) Jacob Peyron <jacob@peyron.io>
// This code is licensed under MIT license.
package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

const (
	NotImplemented = "[not implemented]"
)

// funcMap supplies the map of functions that are available to the renderer.
func funcMap() template.FuncMap {
	sprigMap := sprig.TxtFuncMap()

	// tmpl functions
	tmplMap := template.FuncMap{
		"run": run,

		// Helm compatibility
		"toYaml":        toYaml,
		"fromYaml":      fromYaml,
		"fromYamlArray": fromYamlArray,
		"toJson":        toJson,
		"fromJson":      fromJson,
		"fromJsonArray": fromJsonArray,
		"include":       notImpl2,
		"tpl":           notImpl2,
		"required":      notImpl2e,
		"lookup":        notImpl4e,
	}

	maps.Copy(tmplMap, sprigMap)
	return tmplMap
}

// run arbitrary commands in your template and get the combined output (stdout + stderr)
func run(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		// Should we really handle errors like this?
		if len(out) > 0 {
			return fmt.Sprintf("%v; error: %v", trim(string(out)), err.Error())
		}
		return fmt.Sprintf("error: %v", err.Error())
	}
	return trim(string(out))
}

// toYaml marshals YAML and returns it as a string. Any errors are ignored (Helm compatibility).
func toYaml(str any) string {
	out, err := yaml.Marshal(str)
	if err != nil {
		return ""
	}
	return trim(string(out))
}

// fromYaml unmarshals YAML and returns map[string]any. Any errors are ignored (Helm compatibility).
func fromYaml(str string) map[string]any {
	m := map[string]any{}
	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// fromYamlArray unmarshals a yaml array into []any
func fromYamlArray(str string) []any {
	a := []any{}
	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		a = []any{err.Error()}
	}
	return a
}

// toJson marshals JSON and returns it as a string. Any errors are ignored (Helm compatibility).
func toJson(str string) string {
	out, err := json.Marshal(str)
	if err != nil {
		return ""
	}
	return trim(string(out))
}

// fromJson unmarshals JSON and returns map[string]any. Any errors are ignored (Helm compatibility).
func fromJson(str string) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// fromJsonArray unmarshals a JSON array into []any
func fromJsonArray(str string) []any {
	a := []any{}
	if err := json.Unmarshal([]byte(str), &a); err != nil {
		a = []any{err.Error()}
	}
	return a
}

// trim input message to remove trailing newlines
func trim(str string) string {
	return strings.TrimSuffix(str, "\n")
}

func notImpl2(any, any) string                     { return NotImplemented }
func notImpl2e(any, any) (string, error)           { return NotImplemented, nil }
func notImpl4e(any, any, any, any) (string, error) { return NotImplemented, nil }
