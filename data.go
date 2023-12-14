//	  _    __  __ ,___  _
//	 | |  /  \/  \  _ \| |
//	<   > | |\/| | |_) | |
//	 | |__| |  | |  __/| |__
//	 \___/|_|  |_|_|   \___/
//
// Copyright (C) Jacob Peyron <jacob@peyron.io>
// This code is licensed under MIT license.
package main

// MetaData is the meta data supplied to the renderer
type MetaData struct {
	Name string // Output name (full render pass name)
	Environment
	Values
}

// Document represents some content that should be rendered
type Document struct {
	Filename string
	Content  string
	silent   bool
}

// DocData is data made available to template documents
type DocData struct {
	Filename string
	Content  string
	Env      Environment
	Values
}

type Environment = map[string]string
type Values = map[string]interface{}

func NewMeta(name string, env Environment, vals Values) MetaData {
	return MetaData{
		Name:        name,
		Environment: env,
		Values:      vals,
	}
}

func NewDoc(filename, content string, silent bool) Document {
	return Document{
		Filename: filename,
		Content:  content,
		silent:   silent,
	}
}

func NewDocData(meta MetaData, doc Document) DocData {
	return DocData{
		Filename: doc.Filename,
		Content:  doc.Content,
		Env:      meta.Environment,
		Values:   meta.Values,
	}
}
