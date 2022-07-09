package env_test

import (
	"fmt"
	"net/http"
	"strconv"

	"go.strv.io/env"
)

type debug bool

func (d *debug) UnmarshalText(text []byte) error {
	if d == nil {
		return fmt.Errorf("debug: UnmarshalText: nil pointer")
	}

	b, err := strconv.ParseBool(string(text))
	if err != nil {
		return err
	}
	*d = debug(b)
	return nil
}

type serviceConfig struct {
	Addr    string        `env:"ADDR"`
	Debug   *debug        `env:"DEBUG"`
	Metrics metricsConfig `env:",dive"`
}

type metricsConfig struct {
	Namespace string `env:"METRICS_NAMESPACE"`
}

func ExampleApply() {
	cfg := serviceConfig{}
	env.MustApply(&cfg)

	fmt.Println("Starting HTTP server on address: ", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, nil); err != nil {
		panic(err)
	}
}
