package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type metrics struct {
	ScrapePeriod int64  `env:"METRICS_SCRAPE_PERIOD"`
	Prefix       string `env:"METRICS_PREFIX"`
}

type testConfig struct {
	ignored          struct{} // Unexported fields are skipped. Adding env to unexported field is pointless.
	LogLevel         string   `env:"LOG_LEVEL"`
	RetentionDays    uint32   `env:"RETENTION_DAYS"`
	Scale            float64  `env:"SCALE"`
	BufferingEnabled bool     `env:"BUFFERING_ENABLED"`
	Metrics          metrics
}

func TestFillConfigFromEnv(t *testing.T) {
	testCases := []struct {
		name                   string
		expectedCfg            testConfig
		inputCfg               testConfig
		envLogLevel            string
		envMetricsPrefix       string
		envMetricsScrapePeriod string
		envRetentionDays       string
		envScale               string
		envBufferingEnabled    string
		shouldErr              bool
	}{
		{
			name: "fill-everything/empty-config",
			expectedCfg: testConfig{
				LogLevel:         "debug",
				RetentionDays:    7,
				Scale:            0.8,
				BufferingEnabled: true,
				Metrics: metrics{
					ScrapePeriod: 8,
					Prefix:       "dev",
				},
			},
			inputCfg:               testConfig{},
			envLogLevel:            `debug`,
			envMetricsPrefix:       `dev`,
			envMetricsScrapePeriod: `8`,
			envRetentionDays:       `7`,
			envScale:               `0.8`,
			envBufferingEnabled:    `true`,
		},
		{
			name: "fill-some-fields/empty-config",
			expectedCfg: testConfig{
				LogLevel: "",
				Metrics: metrics{
					ScrapePeriod: 15,
					Prefix:       "experiment",
				},
			},
			inputCfg:               testConfig{},
			envMetricsPrefix:       `experiment`,
			envMetricsScrapePeriod: `15`,
		},
		{
			name: "fill-some-fields/override-config-1",
			expectedCfg: testConfig{
				LogLevel: "info",
				Metrics: metrics{
					ScrapePeriod: 6,
					Prefix:       "stg",
				},
			},
			inputCfg: testConfig{
				LogLevel: "debug",
				Metrics: metrics{
					ScrapePeriod: 10,
					Prefix:       "prod",
				},
			},
			envLogLevel:            `info`,
			envMetricsPrefix:       `stg`,
			envMetricsScrapePeriod: `6`,
		},
		{
			name: "fill-some-fields/override-config-2",
			expectedCfg: testConfig{
				LogLevel: "debug",
				Metrics: metrics{
					ScrapePeriod: 6,
					Prefix:       "stg",
				},
			},
			inputCfg: testConfig{
				LogLevel: "debug",
			},
			envMetricsPrefix:       `stg`,
			envMetricsScrapePeriod: `6`,
		},
	}

	appPrefix := "TEST"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Need to unset env because it persists between test cases.
			if tc.envLogLevel != "" {
				os.Setenv(appPrefix+"_LOG_LEVEL", tc.envLogLevel)
			} else {
				os.Unsetenv(appPrefix + "_LOG_LEVEL")
			}
			if tc.envMetricsPrefix != "" {
				os.Setenv(appPrefix+"_METRICS_PREFIX", tc.envMetricsPrefix)
			} else {
				os.Unsetenv(appPrefix + "_METRICS_PREFIX")
			}
			if tc.envMetricsScrapePeriod != "" {
				os.Setenv(appPrefix+"_METRICS_SCRAPE_PERIOD", tc.envMetricsScrapePeriod)
			} else {
				os.Unsetenv(appPrefix + "_METRICS_SCRAPE_PERIOD")
			}
			if tc.envRetentionDays != "" {
				os.Setenv(appPrefix+"_RETENTION_DAYS", tc.envRetentionDays)
			} else {
				os.Unsetenv(appPrefix + "_RETENTION_DAYS")
			}
			if tc.envScale != "" {
				os.Setenv(appPrefix+"_SCALE", tc.envScale)
			} else {
				os.Unsetenv(appPrefix + "_SCALE")
			}
			if tc.envBufferingEnabled != "" {
				os.Setenv(appPrefix+"_BUFFERING_ENABLED", tc.envBufferingEnabled)
			} else {
				os.Unsetenv(appPrefix + "_BUFFERING_ENABLED")
			}

			cfg := tc.inputCfg
			err := FillConfigFromEnv(&cfg, appPrefix)

			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedCfg, cfg)
			}
		})
	}
}
