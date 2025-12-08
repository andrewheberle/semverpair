# semverpair

[![Go Report Card](https://goreportcard.com/badge/github.com/andrewheberle/semverpair?logo=go&style=flat-square)](https://goreportcard.com/report/github.com/andrewheberle/semverpair)
[![codecov](https://codecov.io/gh/andrewheberle/semverpair/branch/main/graph/badge.svg?token=4KkBld5tkj)](https://codecov.io/gh/andrewheberle/semverpair)

This is a simple CLI tool using [github.com/bep/semverpair](github.com/bep/semverpair) to
encode two versions as a single semver compatible version.

## Installation

```sh
go install github.com/andrewheberle/semverpair/cmd/semverpair
```

## Usage

The CLI will encode a pair of versions as follows:

```sh
semverpair encode [--first|-f] FIRSTVERSION [--second|-s] SECONDVERSION
```

Decoding of a version pair can be performed as follows:

```sh
semverpair decode [--version|-v] VERSION
```
