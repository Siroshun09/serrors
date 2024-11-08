# serrors

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/Siroshun09/serrors)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Siroshun09/serrors/ci.yml?branch=main)
![GitHub](https://img.shields.io/github/license/Siroshun09/serrors)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Siroshun09/serrors)

A Go library to create errors with stack traces.

## Requirements

- Go 1.22+

## Installation

```shell
go get github.com/Siroshun09/serrors
```

## Usage

### Creating errors

- With message: `serrors.New("msg")`
- With format and args: `serrors.Errorf("msg: %s", "hello")`
- With another error: `serrors.WithStackTrace(err)`

### Getting StackTrace attached to errors

- `serrors.GetStackTrace(err)`
- `serrors.GetAttachedStackTrace(err)` 

## License

This project is under the Apache License version 2.0. Please see LICENSE for more info.

Copyright © 2024, Siroshun09
