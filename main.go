package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

var (
	// valuesIn   string
	templateIn []string
	data       Data
)

type Data struct {
	Env    map[string]string
	Values map[string]interface{}
}

func init() {
	// TODO: Init sprig funcs

	// Parse parameters
	// flag.StringVar(&valuesIn, "d", "", "YAML data file")
	flag.Parse()
	templateIn = flag.Args() // positional params (template inputs)

	initEnv()
}

func initEnv() {
	env := os.Environ()
	data.Env = make(map[string]string, len(env))
	for _, e := range env {
		keyVal := strings.SplitN(e, "=", 2)
		data.Env[keyVal[0]] = keyVal[1]
	}
}

func process(name, in string, out io.Writer) {
	t := template.Must(template.New(name).Parse(in))
	err := t.Execute(out, data)
	if err != nil {
		die(err)
	}
}

func main() {
	for _, f := range templateIn {
		buf, err := os.ReadFile(f)
		if err != nil {
			die(err)
		}
		process(f, string(buf), os.Stdout)
	}

	// If no template file was specified, use stdin
	if len(templateIn) == 0 {
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
