# serrors

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/Siroshun09/serrors)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Siroshun09/serrors/ci.yml?branch=v2/main)
![GitHub](https://img.shields.io/github/license/Siroshun09/serrors)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Siroshun09/serrors)

A Go library to create errors with stack traces and structured attributes.

## Requirements

- Go 1.24+

## Installation

```shell
go get github.com/Siroshun09/serrors/v2
```

## Usage

### Creating errors

```go
// Create a new error with a stack trace
err := serrors.New("something went wrong")

// Create a new error with attributes
err := serrors.New("something went wrong", slog.String("key", "value"))

// Wrap an existing error with a stack trace
err := serrors.Wrap(err)

// Wrap with additional attributes
err := serrors.Wrap(err, slog.Int("code", 42))
```

- `New` and `Wrap` always attach a stack trace.
- If the error already has a stack trace from this package, `Wrap` reuses it instead of adding a new one.
- `Wrap(nil)` returns `nil`.

### Getting stack traces

```go
// Get the attached stack trace (returns false if none is attached)
st, ok := serrors.GetAttachedStackTrace(err)

// Get the attached stack trace, or the current call site's stack trace if none is attached
st := serrors.GetStackTrace(err)

// Get the current call site's stack trace
st := serrors.GetCurrentStackTrace()

// Iterate over all stack traces in an error chain (useful with errors.Join)
for err, st := range serrors.GetStackTraces(err) {
    fmt.Println(err, st)
}
```

### Getting attributes

```go
// Iterate over all slog.Attr values attached to an error chain
for err, attr := range serrors.GetAttrs(err) {
    fmt.Println(err, attr)
}
```

### StackTrace and Frame

`StackTrace` is a `[]Frame`, where each `Frame` holds the function name, file, and line number.

```go
st := serrors.GetCurrentStackTrace()
fmt.Println(st.String()) // prints all frames, one per line

for _, frame := range st {
    fmt.Println(frame.String()) // "pkg.FuncName (file.go:42)"
}
```

Both `Frame` and `StackTrace` implement `encoding.TextAppender` via `AppendText([]byte) ([]byte, error)`.

### Interoperability

Errors created by this package implement `Unwrap() error`, so they work with `errors.Is`, `errors.As`, and `fmt.Errorf("%w", ...)` as usual.

## Example

```go
package main

import (
    "errors"
    "fmt"
    "log/slog"

    "github.com/Siroshun09/serrors/v2"
)

func main() {
    base := errors.New("base error")

    // Wrap with a stack trace and an attribute
    err := serrors.Wrap(base, slog.String("user", "alice"))

    // Retrieve the attached stack trace
    if st, ok := serrors.GetAttachedStackTrace(err); ok {
        fmt.Println("stack trace:")
        fmt.Println(st.String())
    }

    // Iterate over attributes in the error chain
    for e, attr := range serrors.GetAttrs(err) {
        fmt.Printf("error=%v attr=%v\n", e, attr)
    }
}
```

## License

This project is under the Apache License version 2.0. Please see LICENSE for more info.

Copyright © 2024-2026, Siroshun09