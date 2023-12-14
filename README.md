# tmpl

This is a small program to render files using [Go templates](https://pkg.go.dev/text/template).
It takes files (or stdin) as input, an optional values YAML file and outputs *one* rendered
document to stdout. The [functions available](https://masterminds.github.io/sprig) to tmpl comes
from the [sprig module](https://github.com/Masterminds/sprig).

This program has [better](https://github.com/hairyhenderson/gomplate)
[alternatives](https://github.com/belitre/gotpl), please consider using those instead. I wrote
this because I wanted a slim program with few dependencies that follows the Unix philosophy of
"do one thing and one thing well".

**WARNING: This utility should not be used in a production environment, use at your own risk.
This is an immature project. Anything can change, at any moment.**

## Installation

For now, installing is only supported via the Go toolchain:

```bash
go install github.com/bdazl/tmpl@latest
```

## Usage

The most basic use-case is reading from stdin (the available environment variables can be
accessed from`.Env`):

```bash
$ echo "My home: {{ .Env.HOME }}" | tmpl
My home: /home/bdazl
```

The utility is designed to render *one* document, but you may want to combine documents (and
share definitions):

```bash
$ tmpl tmpl/tmpl.t tmpl/info.t
  _    __  __ ,___  _
 | |  /  \/  \  _ \| |
<   > | |\/| | |_) | |
 | |__| |  | |  __/| |__
 \___/|_|  |_|_|   \___/

version: v0.1.0
```

The order of the rendered document is very important, because go templates only lets you make
definitions at the start (the initial document(s)):

```bash
$ tmpl tmpl/info.t tmpl/tmpl.t
error: template: tmpl:1:12: executing "tmpl" at <{{template "tmpl.banner"}}>: template "tmpl.banner" not defined
```

Values can be sourced from a YAML-file, and they exists within the `.Values` scope:

```bash
$ echo "{{ .Values.tmpl.version }}" | tmpl -d tmpl/data.yml
v0.1.0
```
