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
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

const (
	StdinFilename = "-"
)

// ProgramArgs describes the parsed inpu from the program arguments
type ProgramArgs struct {
	dataFile      *string
	fileSepFmt    *string
	templateFiles []string
}

func main() {
	var (
		args = parseProgramArgs()
		env  = environ()
		vals = applyOrDflt(readThenYamlUnmarshal, args.dataFile, nil)
		meta = NewMeta("tmpl", env, vals)
	)

	// If no template file was specified, assume stdin
	if len(args.templateFiles) == 0 {
		args.templateFiles = []string{StdinFilename}
	}

	docs := createDocuments(args.templateFiles, false)

	render(meta, docs, os.Stdout)
}

// render multiple documents into one stream. A document may be silent,
// which means that its content will be rendered to /dev/null
func render(meta MetaData, documents []Document, out io.Writer) {
	var (
		t     = template.New(meta.Name)
		funcs = sprig.TxtFuncMap()
		tf    = t.Funcs(funcs)
	)

	for _, d := range documents {
		parsed, err := tf.Parse(d.Content)
		if err != nil {
			die(err)
		}

		docData := NewDocData(meta, d)
		err = parsed.Execute(out, docData)
		if err != nil {
			die(err)
		}
	}
}

func createDocuments(filenames []string, silent bool) []Document {
	var (
		uniq = rmDupes(filenames)
		docs = make([]Document, len(uniq))
	)
	for i, f := range uniq {
		docs[i] = NewDoc(f, readFile(f), true)
	}
	return docs
}

// parseProgramArgs does what you think it does.
func parseProgramArgs() ProgramArgs {
	args := ProgramArgs{
		dataFile: flag.String("d", "", "Data YAML file"),
	}
	flag.Parse()

	// Parse rest of parameters
	args.templateFiles = flag.Args() // positional params (template inputs)
	return args
}

// environ returns the environment variables as a map.
func environ() Environment {
	env := os.Environ()
	out := make(Environment, len(env))
	for _, e := range env {
		keyVal := strings.SplitN(e, "=", 2)
		out[keyVal[0]] = keyVal[1]
	}
	return out
}

// readThenYamlUnmarshal reads file content and unmarshals it to a general map.
func readThenYamlUnmarshal(filename string) map[string]interface{} {
	if filename == "" {
		return nil
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		die(err)
	}

	var out map[string]interface{}
	err = yaml.Unmarshal(buf, &out)
	if err != nil {
		die(err)
	}
	return out
}

// readFile either reads a file from disk, or consumes from stdin
func readFile(filename string) string {
	var (
		buf []byte
		err error
	)

	if filename == StdinFilename {
		buf, err = io.ReadAll(os.Stdin)
	} else {
		buf, err = os.ReadFile(filename)
	}

	if err != nil {
		die(err)
	}
	return string(buf)
}

// rmDupes preserves order and removes duplicate elements from slice
func rmDupes[T comparable](slice []T) []T {
	out := []T{}
	visit := make(map[T]bool)
	for _, e := range slice {
		if _, val := visit[e]; !val {
			visit[e] = true
			out = append(out, e)
		}
	}
	return out
}

// applyOrDflt returns the output from f as applied with maybe, iff the input is not nil.
// If input is nil then the output will be dflt
func applyOrDflt[I, O any](f func(I) O, maybe *I, dflt O) O {
	if maybe == nil {
		return dflt
	}
	return f(*maybe)
}

// die is called if there was an error and we want to terminate the program
func die(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
