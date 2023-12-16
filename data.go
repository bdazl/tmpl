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
	"strings"

	"slices"
)

type Values = map[string]interface{}

// MetaData is the meta data supplied to the renderer
type MetaData struct {
	Name string // Output name (full render pass name)
	Values
}

// Document represents some content that should be rendered
type Document struct {
	Filename string
	Content  string
	silent   bool
}

func NewMeta(name, subKey string, vals Values) MetaData {
	return MetaData{
		Name:   name,
		Values: subKeyCreate(subKey, vals),
	}
}

func NewDoc(filename, content string, silent bool) Document {
	return Document{
		Filename: filename,
		Content:  content,
		silent:   silent,
	}
}

// subKeyCreate takes a string of sub-keys, separated by a '.', and returns a
// nested map where the first sub-key holds the next and so on, until the last element
// holds the actual data.
func subKeyCreate(subKey string, vals Values) Values {
	if subKey == "" {
		return vals
	}
	prev := vals

	keys := strings.Split(subKey, ".")
	slices.Reverse(keys)

	// Traverse sub-keys in reverse order and consequently wrap the previous map
	// in the next one. In the end you should have a map that looks like:
	// map[first][second]...[last] = vals
	for _, k := range keys {
		next := map[string]any{
			k: prev,
		}
		prev = next
	}
	return prev
}
