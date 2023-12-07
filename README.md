# Go helpers

## Install

```sh
go get github.com/samialdury/go
```

## Usage

### `config` package

- `Parse` parses the environment variables into the given struct annotated with `env` tags. Uses the [env](https://github.com/caarlos0/env) package under the hood.
- `Validate` validates the given struct annotated with `validate` tags and returns translated errors with field names taken from the `env` tags. Uses the [validator](https://github.com/go-playground/validator) package under the hood.

```go
package main

import (
    "fmt"
    "os"

    "github.com/samialdury/go/config"
)

type Config struct {
    Mode string `env:"MODE" envDefault:"prod" validate:"required,oneof=dev prod"`
    Port int    `env:"PORT" envDefault:"3000" validate:"required,gt=1024,lt=65535"`
}

func main() {
    var cfg Config

    // This is just for testing the validation
    os.Setenv("MODE", "test")
    os.Setenv("PORT", "75624")

    if err := config.Parse(&cfg); err != nil {
        // handle error
        fmt.Printf("Error parsing config:\n%v\n", err)
        os.Exit(1)
    }

    if err := config.Validate(&cfg); err != nil {
        // handle error
        fmt.Printf("Error validating config:\n%v\n", err)
        /* Output:
        Error validating config:
        env.MODE must be one of [dev prod]. Got `test`.
        env.PORT must be less than 65,535. Got `75624`.
        */
        os.Exit(1)
    }

    // cfg is now ready to use
}
```

### `validation` package

- `NewValidator` returns a new validator instance along with registered English translator for the validation errors.

## License

[MIT](LICENSE)
