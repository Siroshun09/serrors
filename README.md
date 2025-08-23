# serrors

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/Siroshun09/serrors)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Siroshun09/serrors/ci.yml?branch=main)
![GitHub](https://img.shields.io/github/license/Siroshun09/serrors)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Siroshun09/serrors)

A Go library to create errors with stack traces.

## Requirements

- Go 1.24+

## Installation

```shell
go get github.com/Siroshun09/serrors
```

## Usage

### Creating errors

- With message: `serrors.New("msg")`
- With format and args: `serrors.Errorf("msg: %s", "hello")`
- Wrap an existing error: `serrors.WithStackTrace(err)`
  - If `err` is `nil`, it returns `nil`.
  - If `err` already has a stack trace from this package, it returns `err` as-is.

### Getting stack traces from errors

- `serrors.GetStackTrace(err)` returns a stack trace for `err`.
  - If `err` already has a stack trace attached (created by this package), it returns that.
  - Otherwise, it returns the current call site's stack trace.
  - If `err` is `nil`, it returns `nil`.
- `serrors.GetAttachedStackTrace(err)` returns the attached stack trace and a bool.
  - The bool indicates whether `err` had an attached stack trace.

### Example

```go
package main

import (
    "errors"
    "fmt"

    "github.com/Siroshun09/serrors"
)

func main() {
    base := errors.New("base error")

    // Wrap with stack trace
    err := serrors.WithStackTrace(base)

    // Retrieve the attached stack trace
    st, ok := serrors.GetAttachedStackTrace(err)
    if ok {
        fmt.Println("stack trace attached:")
        fmt.Println(st.String())
    }

    // Or always get a stack trace (attached or current)
    fmt.Println(serrors.GetStackTrace(err))
}
```

### Interoperability

- The wrapped error implements `Unwrap() error`, so it works with `errors.Is` and `errors.As`.
- `fmt.Errorf("...: %w", err)` can be used in combination with these errors as usual.

## License

This project is under the Apache License version 2.0. Please see LICENSE for more info.

Copyright © 2024-2025, Siroshun09
