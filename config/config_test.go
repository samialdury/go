package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	Mode string `env:"MODE" envDefault:"prod" validate:"required,oneof=dev prod"`
	Port int    `env:"PORT" envDefault:"3000" validate:"required,gt=1024,lt=65535"`
}

func clearEnvValues(keys ...string) func() {
	return func() {
		for _, key := range keys {
			os.Unsetenv(key)
		}
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name      string
		setupFunc func()
		expected  Config
		expectErr bool
	}{
		{
			name: "should parse environment variables into a struct with default values",
			expected: Config{
				Mode: "prod",
				Port: 3000,
			},
		},
		{
			name: "should parse environment variables into a struct with custom valid values",
			setupFunc: func() {
				os.Setenv("MODE", "dev")
				os.Setenv("PORT", "4000")
			},
			expected: Config{
				Mode: "dev",
				Port: 4000,
			},
		},
		{
			name: "should fail to parse environment variables into a struct with invalid values",
			setupFunc: func() {
				os.Setenv("MODE", "invalid")
				os.Setenv("PORT", "invalid")
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
				t.Cleanup(clearEnvValues("MODE", "PORT"))
			}

			var cfg Config
			err := Parse(&cfg)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, cfg)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name      string
		setupFunc func()
		expected  Config
		expectErr bool
		errMsg    string
	}{
		{
			name: "should validate environment variables with default values",
			expected: Config{
				Mode: "prod",
				Port: 3000,
			},
		},
		{
			name: "should validate environment variables with custom valid values",
			setupFunc: func() {
				os.Setenv("MODE", "dev")
				os.Setenv("PORT", "4000")
			},
			expected: Config{
				Mode: "dev",
				Port: 4000,
			},
		},
		{
			name: "should fail to validate environment variables with invalid port value",
			setupFunc: func() {
				os.Setenv("PORT", "77777")
			},
			expectErr: true,
			errMsg:    "env.PORT must be less than 65,535. Got `77777`.\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
				t.Cleanup(clearEnvValues("MODE", "PORT"))
			}

			var cfg Config
			err := Parse(&cfg)

			assert.NoError(t, err, "Parse should not return an error")

			err = Validate(&cfg)

			if tc.expectErr {
				assert.Error(t, err, "Validate should return an error")
				assert.Equal(t, tc.errMsg, err.Error(), "Error message mismatch")
			} else {
				assert.NoError(t, err, "Validate should not return an error")
				assert.Equal(t, tc.expected, cfg, "Config struct does not match expected")
			}
		})
	}
}
