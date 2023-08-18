# Zerolog Error Marshaler for cockroachdb/errors

Marshals errors produced by cockroachdb/errors package for log output in zerolog.
Inspired by [github.com/rs/zerolog/pkgerrors](https://github.com/rs/zerolog/tree/master/pkgerrors).

> **Warning**  
> This is not an official product of YUMEMI Inc.


## Usage

### Without Context

```go
package main

import (
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/yumemi-inc/zerolog-cockroachdb-errors"
)

zerolog.ErrorStackMarshaler = cockroachdberrors.MarshalStack
log := zerolog.New(zerolog.NewConsoleWriter())

err := errors.Wrap(errors.New("error message"), "from error")
log.Log().Stack().Err(err).Send()
```

### With Context

```go
package main

import (
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/yumemi-inc/zerolog-cockroachdb-errors"
)

zerolog.ErrorStackMarshaler = cockroachdberrors.MarshalStack
log := zerolog.
	New(zerolog.NewConsoleWriter()).
	With().
	Stack().
	Logger()

err := errors.Wrap(errors.New("error message"), "from error")
log.Log().Err(err).Send()
```

## Output

```json
{
  "error": "from error: error message",
  "stack": [
    {
      "stacktrace": [
        {
          "source": "main.go",
          "line": 10,
          "func": "main"
        }
      ]
    },
    {
      "details": [
        "from error"
      ]
    },
    {
      "stacktrace": [
        {
          "source": "main.go",
          "line": 10,
          "func": "main"
        }
      ]
    },
    {
      "details": [
        "error message"
      ]
    }
  ]
}
```
