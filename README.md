# STRV env

Go Library for runtime environment configuration.

## Description
`env` package leverage Go's reflection functionality for scanning structures. If the `env` tag is found, a value of this tag is used for
env lookup. If the env variable is set, the value in the structure is overridden. There are two ways how to process structures:
- If a tag value contains `,dive`, inner fields of a structure are processed (`Session` in the example).
- If a structure implements the `encoding.TextUnmarshaler` interface, `UnmarshalText` is called for the given structure (`AccessTokenExpiration` or `zap.AtomicLevel` in the example).

A good practice when overriding a config with env variables is to use an app prefix. If the env variable `APP_PREFIX` contains some
value (`MY_APP` for example), each defined env variable has to contain the prefix `MY_APP` (`MY_APP_PORT` in the example). There is also an option
to choose a custom prefix for env variables by calling `MustApplyWithPrefix`.

It may happen that the app needs to consume an env variable set by a third party. It's no exception that a cloud provider sets a port your app needs to listen on
and you are unable to modify it. In this case, there is an option to enhance the `env` tag with the `ignoreprefix` clause (`env:"PORT,ignoreprefix"`). While
other env variables will be searched with a prefix included, `PORT` not.

## Examples
```go
package main

import (
	envx "go.strv.io/env"
	timex "go.strv.io/time"

	"go.uber.org/zap"
)

type config struct {
	Port        uint   `json:"port" env:"PORT"`
	StorageAddr string `json:"storage_addr" env:"STORAGE_ADDR"`
	Session     struct {
		AccessTokenExpiration timex.Duration `json:"access_token_expiration" env:"SESSION_ACCESS_TOKEN_EXPIRATION"`
	} `json:"session" env:",dive"`
	LogLevel zap.AtomicLevel `json:"log_level" env:"LOG_LEVEL"`
}

func main() {
	cfg := config{}
	envx.MustApply(&cfg)
}
```

See detailed [example](./example_test.go).
