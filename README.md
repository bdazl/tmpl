# tmpl

This is a small program to render files using [Go templates](https://pkg.go.dev/text/template). It takes files (or stdin)
as input, an optional values YAML file and outputs rendered documents to stdout.

This program has [better](https://github.com/hairyhenderson/gomplate) [alternatives](https://github.com/belitre/gotpl),
please consider using those instead. I wrote this because I wanted a slim program with few dependencies that follows the
Unix philosophy of "do one thing". I might beef it up in the future and remove that last sentence (Google, I'm looking
at you).

**WARNING: This utility should not be used in a production environment**

## Installation

For now, installing is only supported via the Go toolchain:

```bash
go install github.com/bdazl/tmpl@latest
```

## Usage

The most basic use-case is reading from stdin:

```bash
$ echo "My home: {{ .Env.HOME }}" | tmpl
My home: /home/bdazl
```

You can render templates for files directly:

```bash
$ tmpl file1.txt file2.txt
--- file1.txt ---
content of file1
--- file2.txt ---
content of file2
```

This will output the content of the rendered files with a file separator, that can be specified:

```bash
$ tmpl -s ">>> %v" file1.txt file2.txt
>>> file1.txt
some text

>>> file2.txt
some other text
$ tmpl -s "" file1.txt file2.txt
some text
some other text
```
