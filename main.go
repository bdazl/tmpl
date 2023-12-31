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
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	StdinFilename = "-"
)

// ProgramArgs describes the parsed inpu from the program arguments
type ProgramArgs struct {
	dataFile      *string
	subKey        string
	templateFiles []string
}

func main() {
	var (
		args = parseProgramArgs()
		vals = applyOrDflt(readDataYaml, args.dataFile, nil)
		meta = NewMeta("tmpl", args.subKey, vals)
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
		funcs = funcMap()
		t     = template.New(meta.Name)
		tf    = t.Funcs(funcs)
	)

	for _, d := range documents {
		parsed, err := tf.Parse(d.Content)
		if err != nil {
			die(err)
		}

		// Silent files renders to /dev/null
		var writer io.Writer = out
		if d.silent {
			writer = io.Discard
		}

		err = parsed.Execute(writer, meta.Values)
		if err != nil {
			die(err)
		}
	}
}

// createDocuments reads filenames (including special files like stdin) and reads their contents.
// Silent documents are documents that should not be rendered to stdout
func createDocuments(filenames []string, silent bool) []Document {
	var (
		uniq = rmDupes(filenames)
		docs = make([]Document, len(uniq))
	)
	for i, f := range uniq {
		docs[i] = NewDoc(f, readFile(f), silent)
	}
	return docs
}

// parseProgramArgs does what you think it does.
func parseProgramArgs() ProgramArgs {
	args := ProgramArgs{
		dataFile: flag.String("d", "", "Data YAML file"),
	}
	flag.StringVar(&args.subKey, "r", "", "Root key to place data under")
	flag.Parse()

	// Parse rest of parameters
	args.templateFiles = flag.Args() // positional params (template inputs)
	return args
}

// readDataYaml reads yaml-file content and unmarshals it to a general map.
func readDataYaml(filename string) Values {
	if filename == "" {
		return nil
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		die(err)
	}

	var out Values
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
