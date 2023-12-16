# tmpl

This is a small program to render files using [Go templates](https://pkg.go.dev/text/template).
It takes files (or stdin) as input, an optional values YAML file and outputs *one* rendered
document to stdout. The bulk of the [available functions](https://masterminds.github.io/sprig)
to tmpl comes from the [sprig module](https://github.com/Masterminds/sprig).

This program has a [lot](https://github.com/hairyhenderson/gomplate) of
[alternatives](https://github.com/belitre/gotpl), [most](https://github.com/tmc/tmpl) with the
[same](https://github.com/abcum/tmpl) [name](https://github.com/ukautz/tmpl) (look at me being
original). Please consider looking into those projects. I wrote this because I wanted a slim
program with few dependencies that follows the Unix philosophy of "do one thing and one thing well". 
I also wanted certain behaviour and functions.

**WARNING: This utility should not be used in a production environment, use at your own risk.
This is an immature project. Anything can change, at any moment.**

## Installation

For now, installing is only supported via the Go toolchain:

```bash
go install github.com/bdazl/tmpl@latest
```

## Usage

The most basic use-case is reading from stdin:

```bash
$ echo 'My home: {{ env "HOME" }}' | tmpl
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

## Data

Values can be sourced from a YAML-file. Per default they exists within the root scope (`.`).
In the following example, values are sourced from [tmpl/data.yml](tmpl/data.yml) and placed
in the data root. We use the `tmpl.version` value to render the template:

```bash
$ echo "{{ .tmpl.version }}" | tmpl -d tmpl/data.yml
v0.1.0
```

To define the root, use the `-r` flag:

```bash
$ echo "{{ .Values.tmpl.version }}" | tmpl -d tmpl/data.yml -r Values
v0.1.0
```

## Functions

Most of the [available functions](https://masterminds.github.io/sprig) to tmpl comes from the
[sprig module](https://github.com/Masterminds/sprig)
([implementation details here](https://github.com/Masterminds/sprig/blob/master/functions.go)).

In addition to these functions, [tmpl defines its own](funcs.go):
- `run`: run arbitrary commands and return the combined stdout and stderr (ignores error codes).
- `runErr`: similar to `run`, but if exit code is non-zero (or other error) the error will be printed.
- `exitCode`: run command, but return only its exit code

Some functions mimic [Helm](https://github.com/helm/helm), and these are implemented and works
similarly to how they work in `Helm`:
- `toYaml`, `fromYaml`, `fromYamlArray`, `toJson`, `fromJson`, `fromJsonArray`

The following functions are available, but not yet implemented:
- `include`, `tpl`, `required` and `lookup`
