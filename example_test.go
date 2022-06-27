package env_test

import (
	"encoding/json"
	"fmt"
	"os"

	"go.strv.io/env"
)

type serviceConfig struct {
	Name               string `env:"NAME"`
	RequestTimeoutSecs uint   `env:"REQUEST_TIMEOUT_SECS"`
	ColdStorage        struct {
		Pool string `env:"COLD_STORAGE_POOL"`
		Host string `env:"COLD_STORAGE_HOST"`
		Port string `env:"COLD_STORAGE_PORT"`
	}
	HotStorage struct {
		Host string `env:"HOT_STORAGE_HOST"`
		Port string `env:"HOT_STORAGE_PORT"`
	}
}

// ExampleFillConfigFromEnv setups an environment and runs FillConfigFromEnv to demonstrate the functionality of how
// env variables are loaded into the structure based on the env tags.
func ExampleFillConfigFromEnv() {
	appPrefix := "TEST"
	os.Setenv(appPrefix+"_NAME", `"scatterer"`)
	os.Setenv(appPrefix+"_REQUEST_TIMEOUT_SECS", `6`)

	os.Setenv(appPrefix+"_COLD_STORAGE_POOL", `dev"`)
	os.Setenv(appPrefix+"_COLD_STORAGE_HOST", `localhost`)
	os.Setenv(appPrefix+"_COLD_STORAGE_PORT", `8040`)

	os.Setenv(appPrefix+"_HOT_STORAGE_HOST", `localhost`)
	os.Setenv(appPrefix+"_HOT_STORAGE_PORT", `9040`)

	cfg := serviceConfig{}

	if err := env.FillConfigFromEnv(&cfg, appPrefix); err != nil {
		// Handle error
	}

	// Pretty print
	s, _ := json.MarshalIndent(cfg, "", "\t")
	fmt.Print(string(s))

	// Output:
	//{
	//	"Name": "scatterer",
	//	"RequestTimeoutSecs": 6,
	//	"ColdStorage": {
	//		"Pool": "dev",
	//		"Host": "localhost",
	//		"Port": "8040"
	//	},
	//	"HotStorage": {
	//		"Host": "localhost",
	//		"Port": "9040"
	//	}
	//}
}
