# serrors

A Go library to create errors with stack traces.

## Requirements

- Go 1.23.1+

## Usage

### Creating errors

- With message: `serrors.New("msg")`
- With format and args: `serrors.Errorf("msg: %s", "hello")`
- With another error: `serrors.WithStackTrace(err)`

### Getting StackTrace attached to errors

- `serrors.GetStackTrace(err)`

## License

This project is under the Apache License version 2.0. Please see LICENSE for more info.

Copyright © 2024, Siroshun09
