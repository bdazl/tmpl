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

var (
	fileSepFmt    string
	valuesFile    string
	templateFiles []string
	data          Data
)

// This represents the data available from template documents
type Data struct {
	Env    map[string]string
	Values map[string]interface{}
}

func init() {
	// Parse parameters
	flag.StringVar(&valuesFile, "f", "", "YAML data file")
	flag.StringVar(&fileSepFmt, "s", "--- %v ---", "file separator string, where %v is replaced with filename")
	flag.Parse()
	templateFiles = flag.Args() // positional params (template inputs)

	initEnv()
	initValues()
	initFileSep()
}

func initEnv() {
	env := os.Environ()
	data.Env = make(map[string]string, len(env))
	for _, e := range env {
		keyVal := strings.SplitN(e, "=", 2)
		data.Env[keyVal[0]] = keyVal[1]
	}
}

func initValues() {
	if valuesFile == "" {
		return
	}

	buf, err := os.ReadFile(valuesFile)
	if err != nil {
		die(err)
	}

	err = yaml.Unmarshal(buf, &data.Values)
	if err != nil {
		die(err)
	}
}

func initFileSep() {
	// TODO: We should sanitize inputs from format substrings.
	// Only one %v should be allowed

	// Ensure separator ends with newline
	lastEndl := strings.LastIndex(fileSepFmt, "\n")
	if lastEndl != len(fileSepFmt)-1 {
		fileSepFmt += "\n"
	}
}

func process(name, in string, out io.Writer) {
	t := template.Must(template.New(name).Funcs(sprig.TxtFuncMap()).Parse(in))

	err := t.Execute(out, data)
	if err != nil {
		die(err)
	}
}

func main() {
	var (
		templateCount = len(templateFiles)
		multiFile     = templateCount > 1
	)

	for _, f := range templateFiles {
		buf, err := os.ReadFile(f)
		if err != nil {
			die(err)
		}

		if multiFile {
			fmt.Printf(fileSepFmt, f)
		}

		process(f, string(buf), os.Stdout)
	}

	// If no template file was specified, use stdin
	if templateCount == 0 {
		// Must read until EOF before parsing with template
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			die(err)
		}
		process("stdin", string(buf), os.Stdout)
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
