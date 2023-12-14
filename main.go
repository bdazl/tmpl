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

// App is here instead of a global state
type App struct {
	ProgramArgs
	Data
}

// ProgramArgs describes the parsed inpu from the program arguments
type ProgramArgs struct {
	dataFile      *string
	fileSepFmt    *string
	templateFiles []string
}

// Data is the input data available to
type Data struct {
	Env    map[string]string
	Values map[string]interface{}
}

func main() {
	var (
		args = parseProgramArgs()
		data = Data{
			Env:    environ(),
			Values: applyOrDflt(readThenYamlUnmarshal, args.dataFile, nil),
		}
		sep           = applyOrNil(sanitizeFormat, args.fileSepFmt)
		templateCount = len(args.templateFiles)
		multiFile     = templateCount > 1
	)

	for _, f := range args.templateFiles {
		buf, err := os.ReadFile(f)
		if err != nil {
			die(err)
		}

		if multiFile && sep != nil {
			fmt.Printf(*sep, f)
		}

		render(f, string(buf), data, os.Stdout)
	}

	// If no template file was specified, use stdin
	if templateCount == 0 {
		// Must read until EOF before parsing with template
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			die(err)
		}
		render("stdin", string(buf), data, os.Stdout)
	}
}

// render takes a name (used for errors) a (go template) document.
// The document is parsed and executed with the data context into the out writer
func render(name, document string, data Data, out io.Writer) {
	var (
		t     = template.New(name)
		funcs = sprig.TxtFuncMap()
	)

	parsed, err := t.Funcs(funcs).Parse(document)
	if err != nil {
		die(err)
	}

	err = parsed.Execute(out, data)
	if err != nil {
		die(err)
	}
}

// parseProgramArgs does what you think it does.
func parseProgramArgs() ProgramArgs {
	const (
		sepHelp  = "file separator string, where %v is replaced with filename"
		sepDeflt = "--- %v ---"
	)

	args := ProgramArgs{
		dataFile:   flag.String("f", "", "YAML data file"),
		fileSepFmt: flag.String("s", sepDeflt, sepHelp),
	}
	flag.Parse()

	// Parse rest of parameters
	args.templateFiles = flag.Args() // positional params (template inputs)
	return args
}

// environ returns the environment variables as a map.
func environ() map[string]string {
	env := os.Environ()
	out := make(map[string]string, len(env))
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

// sanitizeFormat is already a legacy function and will soon be removed.
func sanitizeFormat(format string) string {
	// TODO: We should not use a Printf format, but go templates
	// Ensure separator ends with newline
	lastEndl := strings.LastIndex(format, "\n")
	if lastEndl != len(format)-1 {
		return format + "\n"
	}
	return format
}

// applyOrDflt returns the output from f as applied with maybe, iff the input is not nil.
// If input is nil then the output will be dflt
func applyOrDflt[I, O any](f func(I) O, maybe *I, dflt O) O {
	if maybe == nil {
		return dflt
	}
	return f(*maybe)
}

// applyOrDflt returns a pointer of output from f as applied with maybe, iff the input is not nil.
// If input is nil then the output will be nil
func applyOrNil[I, O any](f func(I) O, maybe *I) *O {
	if maybe == nil {
		return nil
	}
	o := f(*maybe)
	return &o
}

// die is called if there was an error and we want to terminate the program
func die(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
