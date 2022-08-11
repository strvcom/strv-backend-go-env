package env_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	envx "go.strv.io/env"
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
	err := os.Setenv("APP_PREFIX", "EXAMPLE")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("EXAMPLE_ADDR", ":8080")
	if err != nil {
		panic(err)
	}

	envx.MustApply(&cfg)

	//nolint:gosec
	server := &http.Server{Addr: cfg.Addr}
	fmt.Println("Starting HTTP server on address: ", cfg.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	// Output: Starting HTTP server on address:  :8080
}
