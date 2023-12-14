# tmpl

This is a small program to render files using [Go templates](https://pkg.go.dev/text/template). It takes files (or stdin)
as input, an optional values YAML file and outputs *one* rendered document to stdout.

This program has [better](https://github.com/hairyhenderson/gomplate) [alternatives](https://github.com/belitre/gotpl),
please consider using those instead. I wrote this because I wanted a slim program with few dependencies that follows the
Unix philosophy of "do one thing and one thing well".

**WARNING: This utility should not be used in a production environment, use at your own risk. This is an immature
project. Anything can change, at any moment.**

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

The utility is designed to render *one* document, but you may want to combine documents (and share definitions):

```bash
$ tmpl tmpl/tmpl.t tmpl/info.t
  _    __  __ ,___  _
 | |  /  \/  \  _ \| |
<   > | |\/| | |_) | |
 | |__| |  | |  __/| |__
 \___/|_|  |_|_|   \___/

version: v0.1.0
```
